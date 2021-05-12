package receipts

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"

	// load the postgres river
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"groceryspend.io/server/utils"
)

// AggregatedCategory An aggregation of spend by category
type AggregatedCategory struct {
	Category string
	Value    float32
}

// ReceiptRepository contains the common storage/access patterns for receipts
type ReceiptRepository interface {
	SaveReceipt(receipt *ParsedReceipt) error
	SaveReceiptRequest(request *UnparsedReceiptRequest) error
	GetReceipts(user uuid.UUID) ([]*ParsedReceipt, error)
	GetReceiptDetail(userID uuid.UUID, receiptID uuid.UUID) (*ParsedReceipt, error)
	AggregateSpendByCategoryOverTime(user uuid.UUID, start time.Time, end time.Time) ([]*AggregatedCategory, error)
}

// PostgresReceiptRepository is an implementation of the receipt datastore using postgres
type PostgresReceiptRepository struct {
	DbConnection *sqlx.DB
}

// NewPostgresReceiptRepository creates a new PostgresUserRepo
func NewPostgresReceiptRepository() *PostgresReceiptRepository {
	dbConn, err := sqlx.Open("postgres", utils.GetOsValue("RECEIPTS_POSTGRES_CONN_STR"))

	if err != nil {
		panic("failed to connect to postgres db for users")
	}

	retval := PostgresReceiptRepository{DbConnection: dbConn}
	return &retval
}

// SaveReceipt store parsed receipt to database
func (r *PostgresReceiptRepository) SaveReceipt(receipt *ParsedReceipt) error {

	tx, err := r.DbConnection.BeginTx(context.Background(), &sql.TxOptions{Isolation: 0})
	if err != nil {
		return err
	}

	sql := `
	INSERT INTO parsed_receipts (
		order_number, order_timestamp, sales_tax, tip, service_fee, delivery_fee, discounts, unparsed_receipt_request_id
	)
	VALUES( $1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (order_number) DO UPDATE SET
		order_number = EXCLUDED.order_number,
		order_timestamp = EXCLUDED.order_timestamp,
		sales_tax = EXCLUDED.sales_tax, 
		tip = EXCLUDED.tip, 
		service_fee = EXCLUDED.service_fee,
		delivery_fee = EXCLUDED.delivery_fee, 
		discounts = EXCLUDED.discounts, 
		unparsed_receipt_request_id = EXCLUDED.unparsed_receipt_request_id
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
		receipt.UnparsedReceiptRequestID)
	var prID uuid.UUID
	prRS.Scan(&prID)
	receipt.ID = prID

	// now go through each parsed item and save those, noting the parsed
	for _, pi := range receipt.ParsedItems {
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

func saveParsedItem(tx *sql.Tx, prID uuid.UUID, pi *ParsedItem) error {
	sql := `
	INSERT INTO parsed_items (
		name, total_cost, parsed_receipt_id, category, unit_cost, qty, weight, container_size, container_unit
	)
	VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`
	piRS := tx.QueryRowContext(context.Background(), sql,
		pi.Name,
		pi.TotalCost,
		prID,
		pi.Category,
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
func (r *PostgresReceiptRepository) SaveReceiptRequest(request *UnparsedReceiptRequest) error {
	sql := `
	INSERT INTO unparsed_receipt_requests (
		user_id, original_url, request_timestamp, raw_html
	)
	VALUES( $1, $2, $3, $4)
	ON CONFLICT (original_url) DO UPDATE SET
		user_id = EXCLUDED.user_id,
		original_url = EXCLUDED.original_url,
		request_timestamp = EXCLUDED.request_timestamp, 
		raw_html = EXCLUDED.raw_html
	RETURNING id
	`
	urr := r.DbConnection.QueryRowContext(context.Background(), sql,
		request.UserUUID,
		request.OriginalURL,
		request.RequestTimestamp,
		request.RawHTML,
	)
	var urrID uuid.UUID
	err := urr.Scan(&urrID)
	if err != nil {
		return err
	}
	request.ID = urrID

	return nil
}

// GetReceipts return all receipts for the given user
func (r *PostgresReceiptRepository) GetReceipts(userID uuid.UUID) ([]*ParsedReceipt, error) {
	retval := []*ParsedReceipt{}
	sql := `
		SELECT
			pr.*
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
		var tmp ParsedReceipt
		rows.StructScan(&tmp)
		retval = append(retval, &tmp)
	}

	return retval, nil

}

// GetReceiptDetail return specific receipt details
func (r *PostgresReceiptRepository) GetReceiptDetail(userID uuid.UUID, receiptID uuid.UUID) (*ParsedReceipt, error) {
	retval := ParsedReceipt{}

	// two queries - #1, get the parsed receipt
	sql := `
		SELECT
			pr.*
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
			pi.*
		FROM parsed_items pi
		WHERE pi.parsed_receipt_id = $1
	`
	rows, err := r.DbConnection.QueryxContext(context.Background(), sql, receiptID)
	if err != nil {
		return &retval, err
	}
	defer rows.Close()

	items := []*ParsedItem{}
	for rows.Next() {
		tmp := ParsedItem{}
		rows.StructScan(&tmp)
		items = append(items, &tmp)
	}

	// add items to the parsed receipt
	retval.ParsedItems = items
	return &retval, nil
}

// AggregateSpendByCategoryOverTime get spend by category over time
func (r *PostgresReceiptRepository) AggregateSpendByCategoryOverTime(userID uuid.UUID, start time.Time, end time.Time) ([]*AggregatedCategory, error) {
	sql := `
		select sum(total_cost) as value, category
		from parsed_items pi
		inner join parsed_receipts pr on
			pi.parsed_receipt_id = pr.id
		inner join unparsed_receipt_requests urr on
			pr.unparsed_receipt_request_id = urr.id
		where urr.user_id = $1 
			AND pr.order_timestamp between $2 AND $3
		group by category
		order by sum(total_cost) desc
	`
	retval := []*AggregatedCategory{}
	rows, err := r.DbConnection.QueryxContext(context.Background(), sql, userID, start, end)
	if rows == nil {
		return retval, fmt.Errorf("Got null rows back")
	}
	defer rows.Close()

	if err != nil {
		return retval, err
	}

	for rows.Next() {
		var catSum AggregatedCategory
		err = rows.StructScan(&catSum)
		if err != nil {
			return retval, err
		}
		retval = append(retval, &catSum)
	}

	println(fmt.Sprintf("Num of rows returned: %v", len(retval)))
	return retval, nil
}

type InMemoryReceiptRepository struct {
	idToRequest map[string]*UnparsedReceiptRequest
	idToReceipt map[string]*ParsedReceipt
}

func (r *InMemoryReceiptRepository) SaveReceipt(receipt *ParsedReceipt) error {
	r.idToReceipt[receipt.ID.String()] = receipt
	return nil
}

func (r *InMemoryReceiptRepository) SaveReceiptRequest(request *UnparsedReceiptRequest) error {
	r.idToRequest[request.ID.String()] = request
	return nil
}

func (r *InMemoryReceiptRepository) GetReceipts(user uuid.UUID) ([]*ParsedReceipt, error) {
	retval := []*ParsedReceipt{}

	for _, value := range r.idToRequest {
		if value.UserUUID.String() == user.String() {
			retval = append(retval, value.ParsedReceipt)
		}
	}

	return retval, nil
}

func (r *InMemoryReceiptRepository) GetReceiptDetail(userID uuid.UUID, receiptID uuid.UUID) (*ParsedReceipt, error) {
	// we ignore user id for simpilicity and this was only meant for testing
	return r.idToReceipt[receiptID.String()], nil
}

func (r *InMemoryReceiptRepository) AggregateSpendByCategoryOverTime(user uuid.UUID, start time.Time, end time.Time) ([]*AggregatedCategory, error) {
	return nil, fmt.Errorf("not implemented")
}
