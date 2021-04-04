package receipts

import (
	"fmt"
	"strconv"
	"strings"
)

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

type UnparsedReceiptRequest struct {
	OriginalUrl  string
	IsoTimestamp string
	RawHtml      string
}

type ParsedContainerSize struct {
	Size float32
	Unit string // would prefer type-safe units, but we cannot guarnatee they'll be accurate
}

type ParsedItem struct {
	UnitCost      float32
	Qty           int
	Weight        float32
	TotalCost     float32
	ContainerSize ParsedContainerSize
	Name          string
}

func (p ParsedItem) String() string {
	return fmt.Sprintf("%v: %v", p.Name, p.TotalCost)
}

type ParsedReceipt struct {
	OriginalUrl string
	OrderNumber string
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
