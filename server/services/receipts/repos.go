package receipts

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"groceryspend.io/server/utils"
)

// AggregatedCategory An aggregation of spend by category
type AggregatedCategory struct {
	Category string
	Value    float32
}

// ReceiptRepository contains the common storage/access patterns for receipts
type ReceiptRepository interface {
	SaveReceipt(receipt *ParsedReceipt) (string, error)
	SaveReceiptRequest(request *UnparsedReceiptRequest) (string, error)
	GetReceipts(user uuid.UUID) ([]*ParsedReceipt, error)
	GetReceiptDetail(receiptID uuid.UUID) (*ParsedReceipt, error)
	AggregateSpendByCategoryOverTime(user uuid.UUID, start time.Time, end time.Time) ([]*AggregatedCategory, error)
}

// PostgresReceiptRepository is an implementation of the receipt datastore using postgres
type PostgresReceiptRepository struct {
	DbConnection *gorm.DB
}

// NewPostgresReceiptRepository creates a new PostgresUserRepo
func NewPostgresReceiptRepository() *PostgresReceiptRepository {
	dbConn, err := gorm.Open(postgres.Open(utils.GetOsValue("RECEIPTS_POSTGRES_CONN_STR")), &gorm.Config{})

	if err != nil {
		panic("failed to connect to postgres db for users")
	}

	retval := PostgresReceiptRepository{DbConnection: dbConn}

	// TODO: this should be a script that runs as a different user. That way, the user running queries only
	//       has read/write but not create/delete permissions
	dbConn.AutoMigrate(&UnparsedReceiptRequest{}, &ParsedReceipt{}, &ParsedItem{}, &ParsedContainerSize{})

	return &retval
}

// SaveReceipt store parsed receipt to database
func (r *PostgresReceiptRepository) SaveReceipt(receipt *ParsedReceipt) (string, error) {
	r.DbConnection.Save(receipt)
	return receipt.ID.String(), nil

}

// SaveReceiptRequest store the receipt request in the database
func (r *PostgresReceiptRepository) SaveReceiptRequest(request *UnparsedReceiptRequest) (string, error) {
	r.DbConnection.Save(request)
	return request.ID.String(), nil
}

// GetReceipts return all receipts for the given user
func (r *PostgresReceiptRepository) GetReceipts(user uuid.UUID) ([]*ParsedReceipt, error) {
	retval := []*ParsedReceipt{}
	rows, err := r.DbConnection.Find(&[]ParsedReceipt{}).Where("user_id = ?", user).Order("order_timestamp desc").Rows()
	defer rows.Close()

	if err != nil {
		return retval, err
	}

	for rows.Next() {
		var tmp ParsedReceipt
		r.DbConnection.ScanRows(rows, &tmp)
		retval = append(retval, &tmp)
	}

	return retval, nil

}

// GetReceiptDetail return specific receipt details
func (r *PostgresReceiptRepository) GetReceiptDetail(receiptID uuid.UUID) (*ParsedReceipt, error) {
	var retval *ParsedReceipt
	r.DbConnection.Preload(clause.Associations).First(&retval, receiptID)
	return retval, nil
}

// AggregateSpendByCategoryOverTime get spend by category over time
func (r *PostgresReceiptRepository) AggregateSpendByCategoryOverTime(user uuid.UUID, start time.Time, end time.Time) ([]*AggregatedCategory, error) {
	sql := `
		select sum(total_cost) as value, category
		from parsed_items pi
		inner join parsed_receipts pr on
			pi.parsed_receipt_id = pr.id
		where pr.order_timestamp between ? and ?
		group by category
		order by sum(total_cost) desc
	`
	retval := []*AggregatedCategory{}
	rows, err := r.DbConnection.Raw(sql, start, end).Rows()
	defer rows.Close()

	if err != nil {
		return retval, err
	}

	for rows.Next() {
		var catSum AggregatedCategory
		r.DbConnection.ScanRows(rows, &catSum)
		retval = append(retval, &catSum)
	}

	return retval, nil
}
