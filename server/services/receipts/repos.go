package receipts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/google/uuid"

	// load the postgres river
	_ "github.com/lib/pq"

	// load source file driver
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"

	"groceryspend.io/server/services/categorize"
	"groceryspend.io/server/utils"
)

// ############################## //
// ##        WARNING           ## //
// ## Update this to match the ## //
// ## desired database version ## //
// ## for this git commit      ## //
// ############################## //

// A common question is "why not just always use latest"? The reasoning
// is to allow an outside process to upgrade the database first before rolling out the
// code that works for that specific schema version

// DatabaseVersion is the desired database version for this git commit
const DatabaseVersion = 8

// ReceiptRepository contains the common storage/access patterns for receipts
type ReceiptRepository interface {
	SaveReceipt(receipt *ReceiptDetail) error
	SaveReceiptRequest(request *ParseReceiptRequest) error
	PatchReceiptRequest(request *ParseReceiptRequest) error
	GetReceipts(user uuid.UUID) ([]*ReceiptSummary, error)
	GetReceiptDetail(userID uuid.UUID, receiptID uuid.UUID) (*ReceiptDetail, error)
	PatchReceiptItem(userID uuid.UUID, receiptID uuid.UUID, itemID uuid.UUID, req PatchReceiptItem) error
	AggregateSpendByCategoryOverTime(user uuid.UUID, start time.Time, end time.Time) ([]*AggregatedCategory, error)
}

// DefaultReceiptRepository is an implementation of the receipt datastore using postgres and rabbitmq
type DefaultReceiptRepository struct {
	CatClient       categorize.Client
	DbConnection    *sqlx.DB
	RabbitMqConn    *amqp.Connection // do i need this?
	RabbitMqChannel *amqp.Channel
	RabbitMqQueue   *amqp.Queue // maybe I just need the name?
	RabbitMqDLQ     *amqp.Queue // maybe I just need the name?
}

// NewDefaultReceiptRepository creates a new PostgresUserRepo
func NewDefaultReceiptRepository() *DefaultReceiptRepository {
	dbConn, err := sqlx.Open("postgres", utils.GetOsValue("RECEIPTS_POSTGRES_CONN_STR"))

	if err != nil {
		panic("failed to connect to postgres db for receipts")
	}

	// run migration
	migrationPath := "file://./services/receipts/db/migration"
	driver, err := postgres.WithInstance(dbConn.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		log.Fatalf("Unable to get migration instance: %s", err)
	}

	err = m.Migrate(DatabaseVersion)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Database migration failed: %s", err)
	}

	conn, err := amqp.Dial(utils.GetOsValue("RECEIPTS_RABBITMQ_CONN_STR"))
	if err != nil {
		log.Fatalf("failed to connect to rabbit mq: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to connect to open channel: %s", err)
	}

	q, err := ch.QueueDeclare(
		utils.GetOsValue("RECEIPTS_RABBITMQ_WORK_QUEUE"), // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic("failed to declare queue")
	}

	dlq, err := ch.QueueDeclare(
		utils.GetOsValue("RECEIPTS_RABBITMQ_DLQ"), // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic("failed to declare queue")
	}

	retval := DefaultReceiptRepository{
		DbConnection:    dbConn,
		RabbitMqConn:    conn,
		RabbitMqChannel: ch,
		RabbitMqQueue:   &q,
		RabbitMqDLQ:     &dlq,
		CatClient:       categorize.NewDefaultClient(),
	}
	return &retval
}

// SaveReceipt store parsed receipt to database
func (r *DefaultReceiptRepository) SaveReceipt(receipt *ReceiptDetail) error {

	tx, err := r.DbConnection.BeginTx(context.Background(), &sql.TxOptions{Isolation: 0})
	if err != nil {
		return err
	}

	sql := `
	INSERT INTO parsed_receipts (
		order_number, order_timestamp, sales_tax, tip, service_fee, delivery_fee, discounts, unparsed_receipt_request_id, subtotal_cost
	)
	VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (order_number) DO UPDATE SET
		order_number = EXCLUDED.order_number,
		order_timestamp = EXCLUDED.order_timestamp,
		sales_tax = EXCLUDED.sales_tax, 
		tip = EXCLUDED.tip, 
		service_fee = EXCLUDED.service_fee,
		delivery_fee = EXCLUDED.delivery_fee, 
		discounts = EXCLUDED.discounts, 
		unparsed_receipt_request_id = EXCLUDED.unparsed_receipt_request_id,
		subtotal_cost = EXCLUDED.subtotal_cost
	RETURNING id
	`
	prRS := tx.QueryRowContext(context.Background(), sql,
		receipt.OrderNumber,
		receipt.OrderTimestamp,
		receipt.SalesTax,
		receipt.Tip,
		receipt.ServiceFee,
		receipt.DeliveryFee,
		receipt.Discounts,
		receipt.UnparsedReceiptRequestID,
		receipt.SubtotalCost,
	)
	var prID uuid.UUID
	prRS.Scan(&prID)
	receipt.ID = prID

	// now go through each parsed item and save those, noting the parsed
	for _, pi := range receipt.Items {
		err = saveParsedItem(tx, prID, pi)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// finally commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil

}

func saveParsedItem(tx *sql.Tx, prID uuid.UUID, pi *ReceiptItem) error {
	sql := `
	INSERT INTO parsed_items (
		name, total_cost, parsed_receipt_id, category_id, unit_cost, qty, weight, container_size, container_unit
	)
	VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`
	piRS := tx.QueryRowContext(context.Background(), sql,
		pi.Name,
		pi.TotalCost,
		prID,
		pi.Category.ID,
		pi.UnitCost,
		pi.Qty,
		pi.Weight,
		pi.ContainerSize,
		pi.ContainerUnit)
	var piID uuid.UUID
	err := piRS.Scan(&prID)
	if err != nil {
		return err
	}
	pi.ID = piID

	return nil
}

// SaveReceiptRequest store the receipt request in the database
func (r *DefaultReceiptRepository) SaveReceiptRequest(request *ParseReceiptRequest) error {
	sql := `
	INSERT INTO unparsed_receipt_requests (
		user_id, original_url, request_timestamp, raw_html, request_type_id, status_type_id
	)
	VALUES( $1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	urr := r.DbConnection.QueryRowContext(context.Background(), sql,
		request.UserID,
		request.URL,
		request.Timestamp,
		request.Data,
		request.ParseType,
		request.ParseStatus,
	)
	var urrID uuid.UUID
	err := urr.Scan(&urrID)
	if err != nil {
		return err
	}
	request.ID = urrID

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	err = r.RabbitMqChannel.Publish(
		"",                   // exchange
		r.RabbitMqQueue.Name, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return err
	}
	return nil
}

// PatchReceiptRequest - update request row
func (r *DefaultReceiptRepository) PatchReceiptRequest(request *ParseReceiptRequest) error {
	sql := `
	UPDATE unparsed_receipt_requests SET
		original_url = $3, 
		request_timestamp = $4, 
		raw_html = $5, 
		request_type_id = $6, 
		status_type_id = $7
	WHERE id = $1 and user_id = $2
	RETURNING id
	`
	row := r.DbConnection.QueryRowContext(context.Background(), sql,
		request.ID,
		request.UserID,
		request.URL,
		request.Timestamp,
		request.Data,
		request.ParseType,
		request.ParseStatus,
	)

	var urrID uuid.UUID
	err := row.Scan(&urrID)
	if err != nil {
		return err
	}

	return nil
}

// GetReceipts return all receipts for the given user
func (r *DefaultReceiptRepository) GetReceipts(userID uuid.UUID) ([]*ReceiptSummary, error) {
	retval := []*ReceiptSummary{}
	sql := `
		SELECT
			pr.id as ID,
			urr.user_id as UserUUID,
			pr.order_timestamp as OrderTimestamp, 
			urr.original_url as OriginalURL, 
			urr.request_timestamp as RequestTimestamp,
			(
				pr.sales_tax +
				pr.service_fee +
				pr.delivery_fee +
				pr.discounts +
				pr.tip +
				pr.subtotal_cost
			) as TotalCost
		FROM parsed_receipts pr
		INNER JOIN unparsed_receipt_requests urr ON
			pr.unparsed_receipt_request_id = urr.id
		WHERE urr.user_id = $1
		ORDER BY order_timestamp DESC
	`
	rows, err := r.DbConnection.QueryxContext(context.Background(), sql, userID)
	if err != nil {
		return retval, err
	}
	defer rows.Close()

	for rows.Next() {
		var tmp ReceiptSummary
		rows.StructScan(&tmp)
		retval = append(retval, &tmp)
	}

	return retval, nil

}

// GetReceiptDetail return specific receipt details
func (r *DefaultReceiptRepository) GetReceiptDetail(userID uuid.UUID, receiptID uuid.UUID) (*ReceiptDetail, error) {
	retval := ReceiptDetail{}

	// two queries - #1, get the parsed receipt
	sql := `
		SELECT
			pr.ID,
			urr.original_url as OriginalURL,
			urr.request_timestamp as RequestTimestamp,
			pr.order_number as OrderNumber,
			pr.order_timestamp as OrderTimestamp,
			pr.sales_tax as SalesTax,
			pr.tip as Tip,
			pr.service_fee as ServiceFee,
			pr.delivery_fee as DeliveryFee,
			pr.discounts as Discounts
		FROM parsed_receipts pr
		INNER JOIN unparsed_receipt_requests urr ON
			pr.unparsed_receipt_request_id = urr.id
		WHERE urr.user_id = $1 and pr.ID = $2
		ORDER BY order_timestamp DESC
	`
	row := r.DbConnection.QueryRowxContext(context.Background(), sql, userID, receiptID)

	err := row.StructScan(&retval)
	if err != nil {
		return nil, err
	}

	// #2 - get items
	sql = `
		SELECT
			pi.ID,
			pi.unit_cost as UnitCost,
			pi.qty as Qty,
			pi.weight as Weight,
			pi.total_cost as TotalCost,
			pi.name as Name,
			COALESCE(pi.user_category_id, pi.category_id) as CategoryID,
			pi.container_size as ContainerSize,
			pi.container_unit as ContainerUnit
		FROM parsed_items pi
		WHERE pi.parsed_receipt_id = $1
	`
	rows, err := r.DbConnection.QueryxContext(context.Background(), sql, receiptID)
	if err != nil {
		return &retval, err
	}
	defer rows.Close()

	items := []*ReceiptItem{}
	for rows.Next() {
		tmp := ReceiptItem{}
		rows.StructScan(&tmp)

		// fetch category
		cat, err := r.CatClient.GetCategoryByID(uint(tmp.CategoryID))
		if err != nil {
			return nil, err
		}

		tmp.Category = cat
		items = append(items, &tmp)
	}

	// add items to the parsed receipt
	retval.Items = items
	return &retval, nil
}

// PatchReceiptItem updates specific values on an item. Currently only supports category
func (r *DefaultReceiptRepository) PatchReceiptItem(userID uuid.UUID, receiptID uuid.UUID, itemID uuid.UUID, req PatchReceiptItem) error {
	sql := `
		UPDATE parsed_items as pi SET user_category_id = 
			CASE 
				WHEN category_id = $1 THEN null
				ELSE $1
			END
		FROM parsed_receipts pr, unparsed_receipt_requests urr
		WHERE
			pi.parsed_receipt_id = pr.id AND
			pr.unparsed_receipt_request_id = urr.id AND
		  urr.user_id = $2 AND 
			pr.ID = $3 AND
			pi.ID = $4
		RETURNING pi.ID`

	row := r.DbConnection.QueryRowContext(context.Background(), sql,
		req.CategoryID,
		userID,
		receiptID,
		itemID,
	)
	var piID uuid.UUID
	err := row.Scan(&piID)
	if err != nil {
		return err
	}

	return nil
}

// AggregateSpendByCategoryOverTime get spend by category over time
func (r *DefaultReceiptRepository) AggregateSpendByCategoryOverTime(userID uuid.UUID, start time.Time, end time.Time) ([]*AggregatedCategory, error) {
	sql := `
		select sum(total_cost) as value, COALESCE(user_category_id, category_id) as category_id
		from parsed_items pi
		inner join parsed_receipts pr on
			pi.parsed_receipt_id = pr.id
		inner join unparsed_receipt_requests urr on
			pr.unparsed_receipt_request_id = urr.id
		where urr.user_id = $1 
			AND pr.order_timestamp between $2 AND $3
		group by COALESCE(user_category_id, category_id)
		order by sum(total_cost) desc
	`
	retval := []*AggregatedCategory{}
	rows, err := r.DbConnection.QueryxContext(context.Background(), sql, userID, start, end)
	if err != nil {
		return retval, err
	}

	if rows == nil {
		return retval, fmt.Errorf("Got null rows back")
	}
	defer rows.Close()

	for rows.Next() {

		var id int
		var value float32

		err := rows.Scan(&value, &id)
		if err != nil {
			println(err.Error())
			return nil, err
		}

		cat, err := r.CatClient.GetCategoryByID(uint(id))
		if err != nil {
			println(err.Error())
			return nil, err
		}

		retval = append(retval, &AggregatedCategory{
			Category: cat.Name,
			Value:    value,
		})
	}

	return retval, nil
}
