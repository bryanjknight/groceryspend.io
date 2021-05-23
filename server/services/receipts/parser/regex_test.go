package parser

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"groceryspend.io/server/services/receipts"
)

func getTestDataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename)))), "test", "data", "marketbasket")
}

func readFileAsString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func TestRegexParser(t *testing.T) {

	type test struct {
		filename        string
		expectedReceipt *receipts.ReceiptDetail
	}

	tests := []test{
		{
			filename: "receipt1.txt",
			expectedReceipt: &receipts.ReceiptDetail{
				SalesTax:     0.0,
				ServiceFee:   0.0,
				DeliveryFee:  0.0,
				Tip:          0.0,
				Discounts:    0.0,
				SubtotalCost: 34.05,
				Items: []*receipts.ReceiptItem{
					{
						Name:      "GR/HOUSE RED PEPPERS",
						UnitCost:  2.99,
						Weight:    0.55,
						TotalCost: 1.64,
					},
					{
						Name:      "GR/HOUSE RED PEPPERS",
						UnitCost:  2.99,
						Weight:    0.77,
						TotalCost: 1.64,
					},
					{
						Name:      "GR/HOUSE RED PEPPERS",
						UnitCost:  2.99,
						Weight:    0.55,
						TotalCost: 2.30,
					},
					{
						Name:      "RED SEEDLESS GRAPES",
						UnitCost:  1.79,
						Weight:    2.38,
						TotalCost: 1.64,
					},
					{
						Name:      "1# PEELED BABY CARRT",
						UnitCost:  1.29,
						Qty:       1,
						TotalCost: 1.29,
					},
					{
						Name:      "CABOT SERIOUS BRICK",
						UnitCost:  8.99,
						Qty:       1,
						TotalCost: 8.99,
					},
					{
						Name:      "DRAG WHL MOZZ CHUNK",
						UnitCost:  2.50,
						Qty:       1,
						TotalCost: 2.50,
					},
					{
						Name:      "JOE RST RED PEPPERS",
						UnitCost:  3.59,
						Qty:       1,
						TotalCost: 3.59,
					},
					{
						Name:      "FAGE WHOLE 1/2 KILO",
						UnitCost:  2.99,
						Qty:       1,
						TotalCost: 2.99,
					},
					{
						Name:      "GRN MTN SALSA MLD",
						UnitCost:  3.50,
						Qty:       1,
						TotalCost: 3.50,
					},
					{
						Name:      "MB RESTAURANT TORTS",
						UnitCost:  2.00,
						Qty:       1,
						TotalCost: 2.00,
					},
				},
			},
		},
	}

	for _, test := range tests {
		absPath := filepath.Join(getTestDataDir(), test.filename)
		println(absPath)
		text := readFileAsString(absPath)
		actualReceipt, err := RegexParser(text)

		if err != nil {
			t.Fatalf(err.Error())
		}

		if test.expectedReceipt != actualReceipt {
			t.Fatalf("Expected did not match actual receipt")
		}
	}
}
