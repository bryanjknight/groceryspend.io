package receipts

import (
	"fmt"

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
