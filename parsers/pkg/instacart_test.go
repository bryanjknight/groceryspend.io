package pkg

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
	return filepath.Join(filepath.Dir(filepath.Dir(filename)), "test", "data")
}

func readFileAsString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(b)
}

func TestInstacartReceipt(t *testing.T) {
	testDataDir := getTestDataDir()
	orderNumber := "1234567890" // TODO: extract into test table
	fileContent := readFileAsString(filepath.Join(testDataDir, "instacart", fmt.Sprintf("%s.txt", orderNumber)))
	fileContentReader := strings.NewReader(fileContent)

	receiptRequest := UnparsedReceiptRequest{}
	parsedHtml, err := html.Parse(fileContentReader)
	if err != nil {
		t.Errorf("Failed to parse html data: %s", err)
	}
	receiptRequest.Receipt = parsedHtml
	receiptRequest.OriginalUrl = fmt.Sprintf("https://www.instacart.com/orders/%s", orderNumber)

	receipt, err := Parse(receiptRequest)
	if err != nil {
		t.Errorf("Failed to parse receipt: %s", err)
	}

	expectedTotalItems := 27 // TODO: extract into test table
	if len(receipt.ParsedItems) != expectedTotalItems {
		t.Errorf("Expected %v items, got %v", expectedTotalItems, len(receipt.ParsedItems))
	}

	// sum the parsed items to get the subtotal
	expectedSum := float32(150.96) // TODO: extract into test table

	actualSum := float32(0.0)
	for _, item := range receipt.ParsedItems {
		actualSum += item.TotalCost
	}

	if expectedSum != actualSum {
		t.Errorf("Expectd total sum %v, got %v", expectedSum, actualSum)
	}
}
