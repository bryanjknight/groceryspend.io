package receipts

import (
	"errors"

	"golang.org/x/net/html"
)

// ParseWfmHTMLRecipt parse a Whole Foods receipt
func ParseWfmHTMLRecipt(doc *html.Node) (ReceiptDetail, error) {
	return ReceiptDetail{}, errors.New("not implemented")
}
