package receipts

import (
	"fmt"
	"strconv"

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
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OriginalURL  string    `gorm:"notNull"`
	IsoTimestamp string    `gorm:"notNull"`
	RawHTML      string    `gorm:"notNull"`
}

// ParsedContainerSize the size of an item's container (e.g. a 16oz container of strawberries)
type ParsedContainerSize struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Size         float32
	Unit         string
	ParsedItemID uuid.UUID `gorm:"type:uuid,notNull"`
}

// ParsedItem a parsed line item from a receipt
type ParsedItem struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UnitCost        float32
	Qty             int
	Weight          float32
	TotalCost       float32 `gorm:"notNull"`
	ContainerSize   ParsedContainerSize
	Name            string    `gorm:"notNull"`
	ParsedReceiptID uuid.UUID `gorm:"type:uuid;notNull"`
}

func (p ParsedItem) String() string {
	return fmt.Sprintf("%v: %v", p.Name, p.TotalCost)
}

// ParsedReceipt a fully parsed receipt
type ParsedReceipt struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID      string    `gorm:"notNull"`
	OriginalURL string    `gorm:"notNull"`
	OrderNumber string    `gorm:"notNull"`
	Timestamp   string    `gorm:"notNull"`
	ParsedItems []ParsedItem
	// TODO: break out tax, tip, and fees into 1-to-many relationship
	//			 as some jurisdictions could have multiple taxes
	SalesTax    float32
	Tip         float32
	ServiceFee  float32
	DeliveryFee float32
	Discounts   float32
}

func (p ParsedReceipt) String() string {

	builder := strings.Builder{}

	builder.WriteString("Items:\n=====\n")
	for _, item := range p.ParsedItems {
		builder.WriteString(fmt.Sprintf("%v\n", item))
	}

	builder.WriteString("=====\n")
	builder.WriteString(fmt.Sprintf("Delivery Fee: %v\n", p.DeliveryFee))
	builder.WriteString(fmt.Sprintf("Service Fee: %v\n", p.ServiceFee))
	builder.WriteString(fmt.Sprintf("Sales Tax: %v\n", p.SalesTax))
	builder.WriteString(fmt.Sprintf("Tip: %v\n", p.Tip))
	builder.WriteString(fmt.Sprintf("Discounts: %v\n", p.Discounts))

	return builder.String()
}
