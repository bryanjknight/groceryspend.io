package parsers

import (
	"errors"
	"strings"

	"groceryspend.io/server/models"
)

func Parse(request models.UnparsedReceiptRequest) (models.ParsedReceipt, error) {
	if strings.Contains(request.OriginalUrl, "instacart.com") {
		return ParseInstcartHtmlReceipt(request.Receipt)
	}
	if strings.Contains(request.OriginalUrl, "amazon.com") {
		return ParseWfmHtmlRecipt(request.Receipt)
	}

	return models.ParsedReceipt{}, errors.New("unable to match URL with parser")
}
