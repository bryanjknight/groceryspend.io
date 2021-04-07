package receipts

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

// ParseReceipt given a request, try to parse the receipt into something machine readable
func ParseReceipt(request UnparsedReceiptRequest) (ParsedReceipt, error) {

	// parse html
	dataReader := strings.NewReader(request.RawHTML)
	parsedHTML, err := html.Parse(dataReader)
	if err != nil {
		return ParsedReceipt{}, err
	}
	if strings.Contains(request.OriginalURL, "instacart.com") {
		return ParseInstacartHTMLReceipt(parsedHTML)
	}
	if strings.Contains(request.OriginalURL, "amazon.com") {
		return ParseWfmHTMLRecipt(parsedHTML)
	}

	return ParsedReceipt{}, errors.New("unable to match URL with parser")
}
