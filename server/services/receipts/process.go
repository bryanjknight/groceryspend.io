package receipts

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

// ParseReceipt given a request, try to parse the receipt into something machine readable
func ParseReceipt(request ParseReceiptRequest) (ReceiptDetail, error) {

	// parse html
	dataReader := strings.NewReader(request.Data)
	parsedHTML, err := html.Parse(dataReader)
	if err != nil {
		return ReceiptDetail{}, err
	}
	if strings.Contains(request.URL, "instacart.com") {

		receipt, err := ParseInstacartHTMLReceipt(parsedHTML)
		if err != nil {
			return ReceiptDetail{}, err
		}

		// get the order number from the URL
		splitURL := strings.Split(request.URL, "/")
		receipt.OrderNumber = splitURL[len(splitURL)-1]
		return receipt, nil
	}
	if strings.Contains(request.URL, "amazon.com") {
		return ParseWfmHTMLRecipt(parsedHTML)
	}

	return ReceiptDetail{}, errors.New("unable to match URL with parser")
}
