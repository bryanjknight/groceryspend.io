package parsers

import (
	"errors"

	"golang.org/x/net/html"

	"groceryspend.io/data-processing/models"
)

func ParseWfmHtmlRecipt(doc *html.Node) (models.ParsedReceipt, error) {
	return models.ParsedReceipt{}, errors.New("not implemented")
}
