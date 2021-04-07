package receipts

import (
	"errors"

	"golang.org/x/net/html"
)

// ParseWfmHTMLRecipt parse a Whole Foods receipt
func ParseWfmHTMLRecipt(doc *html.Node) (ParsedReceipt, error) {
	return ParsedReceipt{}, errors.New("not implemented")
}
