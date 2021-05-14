package receipts

import "groceryspend.io/server/services/categorize"

// HandleReceiptRequest handles the process of parsing a receipt and saving it
func HandleReceiptRequest(receiptRequest ParseReceiptRequest, repo ReceiptRepository, categorizeClient categorize.Client) error {
	receipt, err := ParseReceipt(receiptRequest)
	if err != nil {
		return err
	}

	// categorize the items
	for _, item := range receipt.Items {
		itemNames := []string{item.Name}
		itemToCat := make(map[string]string)

		err = categorizeClient.GetCategoryForItems(itemNames, &itemToCat)
		if err != nil {
			return err
		}
		item.Category = itemToCat[item.Name]
	}

	receipt.UnparsedReceiptRequestID = receiptRequest.ID

	err = repo.SaveReceipt(&receipt)

	// could be nil or something
	return err
}
