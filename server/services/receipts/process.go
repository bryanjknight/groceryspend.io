package receipts

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

func ParseReceipt(request UnparsedReceiptRequest) (ParsedReceipt, error) {

	// parse html
	dataReader := strings.NewReader(request.RawHtml)
	parsedHtml, err := html.Parse(dataReader)
	if err != nil {
		return ParsedReceipt{}, err
	}
	if strings.Contains(request.OriginalUrl, "instacart.com") {
		return ParseInstacartHtmlReceipt(parsedHtml)
	}
	if strings.Contains(request.OriginalUrl, "amazon.com") {
		return ParseWfmHtmlRecipt(parsedHtml)
	}

	return ParsedReceipt{}, errors.New("unable to match URL with parser")
}
