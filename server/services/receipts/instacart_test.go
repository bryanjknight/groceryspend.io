package receipts

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

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
	}

	tests := []test{
		{OrderNumber: "wegmans-replace-refund", ExpectedTotalItems: 27, ExpectedSubtotal: 150.96, ExpectedTotal: 188.44},
		// Note that bj's has a different subtotal. There's a bug in instacart, so we will calculate subtotal ourselves
		{OrderNumber: "bj-wholesale-all-found", ExpectedTotalItems: 6, ExpectedSubtotal: 202.93, ExpectedTotal: 255.62},
	}

	for _, test := range tests {
		testDataDir := getTestDataDir()
		orderNumber := test.OrderNumber
		fileContent := readFileAsString(filepath.Join(testDataDir, "instacart", fmt.Sprintf("%s.txt", orderNumber)))
		fileContentReader := strings.NewReader(fileContent)

		parsedHtml, err := html.Parse(fileContentReader)
		if err != nil {
			t.Errorf("Failed to parse html data: %s", err)
		}

		receipt, err := ParseInstacartHtmlReceipt(parsedHtml)
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
	}
}
