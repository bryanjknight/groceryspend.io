package pkg

import (
	"errors"

	"golang.org/x/net/html"
)

func ParseWfmHtmlRecipt(doc *html.Node) (ParsedReceipt, error) {
	return ParsedReceipt{}, errors.New("not implemented")
}
