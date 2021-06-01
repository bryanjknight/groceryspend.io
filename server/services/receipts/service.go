package receipts

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"groceryspend.io/server/services/categorize"
)

// HandleReceiptRequest handles the process of parsing a receipt and saving it
func HandleReceiptRequest(
	receiptRequest ParseReceiptRequest,
	repo ReceiptRepository,
	categorizeClient categorize.Client) error {
	receipt, err := ParseReceipt(receiptRequest)
	if err != nil {
		println(fmt.Sprintf("Failed to parse receipt request %s", receiptRequest.ID.String()))
		return err
	}

	// categorize the items
	for _, item := range receipt.Items {
		itemNames := []string{item.Name}
		itemToCat := make(map[string]*categorize.Category)

		err = categorizeClient.GetCategoryForItems(itemNames, &itemToCat)
		if err != nil {
			println(fmt.Sprintf("Failed to get category for %s", item.Name))
			return err
		}
		item.Category = itemToCat[item.Name]
	}

	receipt.UnparsedReceiptRequestID = receiptRequest.ID

	err = repo.SaveReceipt(&receipt)
	if err != nil {
		println(fmt.Sprintf("Failed to save receipt for request %s", receiptRequest.ID.String()))
		return err
	}

	return nil
}

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
