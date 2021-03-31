package pkg

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
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
	OriginalUrl string
	Receipt     *html.Node
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

type ParsedReceipt struct {
	OriginalUrl string
	OrderNumber string
	ParsedItems []ParsedItem
	SalesTax    float32
	Tip         float32
	ServiceFee  float32
	DeliveryFee float32
}
