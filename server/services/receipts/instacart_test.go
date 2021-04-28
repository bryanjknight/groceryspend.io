package receipts

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

// TODO: memoize
func getTestDataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename))), "test", "data")
}

func readFileAsString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func TestInstacartReceipt(t *testing.T) {

	type test struct {
		OrderNumber        string
		ExpectedTotalItems int
		ExpectedSubtotal   float32
		ExpectedTotal      float32
		OrderTimestamp     time.Time
	}

	loc, _ := time.LoadLocation("America/New_York")

	tests := []test{
		{
			OrderNumber:        "wegmans-replace-refund",
			ExpectedTotalItems: 27,
			ExpectedSubtotal:   150.96,
			ExpectedTotal:      188.44,
			// Mar 28, 2021, 9:16 AM
			OrderTimestamp: time.Date(2021, 3, 28, 9, 16, 0, 0, loc),
		},
		// Note that bj's has a different subtotal. There's a bug in instacart, so we will calculate subtotal ourselves
		{
			OrderNumber:        "bj-wholesale-all-found",
			ExpectedTotalItems: 6,
			ExpectedSubtotal:   202.93,
			ExpectedTotal:      255.62,
			// Mar 28, 2021, 10:27 AM
			OrderTimestamp: time.Date(2021, 3, 28, 10, 27, 0, 0, loc),
		},
	}

	for _, test := range tests {
		testDataDir := getTestDataDir()
		orderNumber := test.OrderNumber
		fileContent := readFileAsString(filepath.Join(testDataDir, "instacart", fmt.Sprintf("%s.txt", orderNumber)))
		fileContentReader := strings.NewReader(fileContent)

		parsedHTML, err := html.Parse(fileContentReader)
		if err != nil {
			t.Errorf("Failed to parse html data: %s", err)
		}

		receipt, err := ParseInstacartHTMLReceipt(parsedHTML)
		if err != nil {
			t.Errorf("Failed to parse receipt: %s", err)
		}

		expectedTotalItems := test.ExpectedTotalItems
		if len(receipt.ParsedItems) != expectedTotalItems {
			t.Errorf("Expected %v items, got %v", expectedTotalItems, len(receipt.ParsedItems))
		}

		// sum the parsed items to get the subtotal
		expectedSubtotal := test.ExpectedSubtotal

		actualSubtotal := float32(0.0)
		for _, item := range receipt.ParsedItems {
			actualSubtotal += item.TotalCost
		}

		if expectedSubtotal != actualSubtotal {
			t.Errorf("Expectd subtotal %v, got %v", expectedSubtotal, actualSubtotal)
		}

		expectedTotal := test.ExpectedTotal
		actualTotal := actualSubtotal + receipt.SalesTax + receipt.ServiceFee + receipt.Tip + receipt.DeliveryFee + receipt.Discounts

		if expectedTotal != actualTotal {
			t.Errorf("Expected total %v, got %v", expectedTotal, actualTotal)
		}

		if !test.OrderTimestamp.Equal(receipt.OrderTimestamp) {
			t.Errorf("Expected timestamp %v, got %v", test.OrderTimestamp, receipt.OrderTimestamp)
		}
	}
}
