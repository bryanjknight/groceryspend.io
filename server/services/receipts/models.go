package receipts

import (
	"strconv"
	"time"

	"github.com/google/uuid"

	"strings"
)

// ParseStringToUSDAmount parse a string value into a float32 assuming the value is something like $123.45
func ParseStringToUSDAmount(s string) (float32, error) {
	// strip the $ out, convert to a float32
	sNoDollarSign := strings.ReplaceAll(s, "$", "")
	sTrimmedNoDollarSign := strings.TrimSpace(sNoDollarSign)
	val, err := strconv.ParseFloat(sTrimmedNoDollarSign, 32)
	if err != nil {
		return 0, err
	}
	return float32(val), nil
}

// UnparsedReceiptRequest a request to parse a receipt
type UnparsedReceiptRequest struct {
	ID               uuid.UUID
	UserUUID         uuid.UUID `db:"user_uuid"`
	OriginalURL      string    `db:"original_url"`
	RequestTimestamp time.Time `db:"request_timestamp"`
	RawHTML          string    `db:"raw_html"`
	ParsedReceipt    *ParsedReceipt
}

// ParsedItem a parsed line item from a receipt
type ParsedItem struct {
	ID              uuid.UUID
	UnitCost        float32 `db:"unit_cost"`
	Qty             int
	Weight          float32
	TotalCost       float32 `db:"total_cost"`
	Name            string
	ParsedReceiptID uuid.UUID `db:"parsed_receipt_id"`
	Category        string
	ContainerSize   float32 `db:"container_size"`
	ContainerUnit   string  `db:"container_unit"`
}

// ParsedReceipt a fully parsed receipt
type ParsedReceipt struct {
	ID             uuid.UUID
	OrderNumber    string    `db:"order_number"`
	OrderTimestamp time.Time `db:"order_timestamp"`
	ParsedItems    []*ParsedItem
	// TODO: break out tax, tip, and fees into 1-to-many relationship
	//			 as some jurisdictions could have multiple taxes
	SalesTax                 float32 `db:"sales_tax"`
	Tip                      float32
	ServiceFee               float32 `db:"service_fee"`
	DeliveryFee              float32 `db:"delivery_fee"`
	Discounts                float32
	UnparsedReceiptRequestID uuid.UUID `db:"unparsed_receipt_request_id"`
}
