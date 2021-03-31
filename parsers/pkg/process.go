package pkg

import (
	"errors"
	"strings"
)

func Parse(request UnparsedReceiptRequest) (ParsedReceipt, error) {
	if strings.Contains(request.OriginalUrl, "instacart.com") {
		return ParseInstcartHtml(request.Receipt)
	}

	return ParsedReceipt{}, errors.New("unable to match URL with parser")
}
