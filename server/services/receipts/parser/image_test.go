package parser

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/textract"
	"groceryspend.io/server/services/receipts"
)

func TestTextractResponse(t *testing.T) {

	type test struct {
		file           string
		expectedResult *receipts.ReceiptDetail
		expectedTotal  float32
	}

	tests := []test{
		{
			file:           filepath.Join(getTestDataDir(), "marketbasket", "receipt1-apiResponse.json"),
			expectedTotal:  34.05,
			expectedResult: nil,
		},
		{
			file:           filepath.Join(getTestDataDir(), "hannaford", "receipt1-apiResponse.json"),
			expectedTotal:  29.92,
			expectedResult: nil,
		},
		{
			file:           filepath.Join(getTestDataDir(), "wegmans", "receipt1-apiResponse.json"),
			expectedTotal:  64.01,
			expectedResult: nil,
		},
	}

	for _, testInstance := range tests {
		var response textract.AnalyzeDocumentOutput
		fileText := readFileAsString(testInstance.file)
		reader := strings.NewReader(fileText)
		err := json.NewDecoder(reader).Decode(&response)
		if err != nil {
			println(err.Error())
		}

		receiptDetail, err := ParseImageReceipt(&response, testInstance.expectedTotal)
		if err != nil {
			t.Errorf("error while processing %s: %s", testInstance.file, err.Error())
		} else if receiptDetail == nil {
			t.Errorf("didn't get receipt detail for %s", testInstance.file)
		} else {
			// debug, print the details
			for _, i := range receiptDetail.Items {
				t.Logf("%s: %v", i.Name, i.TotalCost)
			}

		}

	}
}
