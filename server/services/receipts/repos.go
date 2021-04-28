package receipts

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"groceryspend.io/server/utils"
)

// AggregateCategoryResponse An aggregation of spend by category
type AggregateCategoryResponse struct {
	CategoryToSpend map[string]float32
}

// ReceiptRepository contains the common storage/access patterns for receipts
type ReceiptRepository interface {
	AddReceipt(receipt ParsedReceipt) (string, error)
	AddReceiptRequest(request UnparsedReceiptRequest) (string, error)
	AggregateSpendByCategoryOverTime(user uuid.UUID, start time.Time, end time.Time) (*AggregateCategoryResponse, error)
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

// AddReceipt store parsed receipt to database
func (r *PostgresReceiptRepository) AddReceipt(receipt ParsedReceipt) (string, error) {
	r.DbConnection.Create(&receipt)
	return receipt.ID.String(), nil

}

// AddReceiptRequest store the receipt request in the database
func (r *PostgresReceiptRepository) AddReceiptRequest(request UnparsedReceiptRequest) (string, error) {
	r.DbConnection.Create(&request)
	return request.ID.String(), nil
}

// AggregateSpendByCategoryOverTime get spend by category over time
func (r *PostgresReceiptRepository) AggregateSpendByCategoryOverTime(user uuid.UUID, start time.Time, end time.Time) (*AggregateCategoryResponse, error) {

	return nil, nil
}
