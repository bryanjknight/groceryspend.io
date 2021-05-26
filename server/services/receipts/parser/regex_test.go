package parser

import (
	"fmt"
	"path/filepath"
	"testing"

	"groceryspend.io/server/services/receipts"
)

func TestRegexParser(t *testing.T) {

	type test struct {
		filename        string
		expectedReceipt *receipts.ReceiptDetail
	}

	tests := []test{
		{
			filename: "marketbasket/receipt1.txt",
			expectedReceipt: &receipts.ReceiptDetail{
				SalesTax:     0.0,
				ServiceFee:   0.0,
				DeliveryFee:  0.0,
				Tip:          0.0,
				Discounts:    0.0,
				SubtotalCost: 34.05,
				TotalCost:    34.05,
				Items:        []*receipts.ReceiptItem{},
			},
		},
		{
			filename: "marketbasket/receipt2.txt",
			expectedReceipt: &receipts.ReceiptDetail{
				SalesTax:     0.0,
				ServiceFee:   0.0,
				DeliveryFee:  0.0,
				Tip:          0.0,
				Discounts:    0.0,
				SubtotalCost: 117.66,
				TotalCost:    117.66,
				Items:        []*receipts.ReceiptItem{},
			},
		},
		{
			filename: "hannaford/receipt1.txt",
			expectedReceipt: &receipts.ReceiptDetail{
				SalesTax:    0.0,
				ServiceFee:  0.0,
				DeliveryFee: 0.0,
				Tip:         0.0,
				Discounts:   0.0,
				// Hannafords in NH doesn't print a subtotal b/c no tax (Live Free or Die)
				SubtotalCost: 0.0,
				TotalCost:    29.92,
				Items:        []*receipts.ReceiptItem{},
			},
		},
	}

	for _, test := range tests {
		absPath := filepath.Join(getTestDataDir(), test.filename)
		text := readFileAsString(absPath)
		receiptDetail, err := RegexParser(text)

		if err != nil {
			t.Fatalf(err.Error())
		}

		actualTotal := float32(0.0)
		for _, item := range receiptDetail.Items {
			actualTotal += item.TotalCost
		}
		actualTotal += receiptDetail.SalesTax

		// because we store the values as decimals, and float32 and additions is hard
		// we will use string comparison to verify the totals. Ideally I would store everything
		// as cents
		if fmt.Sprintf("%.2f", actualTotal) != fmt.Sprintf("%.2f", test.expectedReceipt.TotalCost) {
			t.Errorf(
				fmt.Sprintf(
					"Expected receipt %s to have total cost of %v, but calculated to be %v",
					test.filename,
					receiptDetail.TotalCost,
					actualTotal,
				),
			)
		}

		if receiptDetail.TotalCost != test.expectedReceipt.TotalCost {
			t.Errorf(
				fmt.Sprintf(
					"Expected total for %s to be %v, got %v",
					test.filename,
					test.expectedReceipt.TotalCost,
					receiptDetail.TotalCost,
				),
			)
		}
	}
}
