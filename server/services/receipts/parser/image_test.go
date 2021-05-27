package parser

import (
	"encoding/json"
	"fmt"
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
		config         *ImageReceiptParseConfig
	}

	tests := []test{
		{
			file: filepath.Join(getTestDataDir(), "marketbasket", "receipt1-apiResponse.json"),
			config: &ImageReceiptParseConfig{
				maxItemDescXPos: 0.7,
				tolerance:       0.015,
			},
			expectedResult: nil,
		},
		{
			file: filepath.Join(getTestDataDir(), "hannaford", "receipt1-apiResponse.json"),
			config: &ImageReceiptParseConfig{
				maxItemDescXPos: 0.7,
				tolerance:       0.015,
			},
			expectedResult: nil,
		},
	}

	for _, testInstance := range tests {
		println("")
		println(fmt.Sprintf("Processing %s", testInstance.file))
		println("")
		var response textract.AnalyzeDocumentOutput
		fileText := readFileAsString(testInstance.file)
		reader := strings.NewReader(fileText)
		err := json.NewDecoder(reader).Decode(&response)
		if err != nil {
			println(err.Error())
		}

		err = ProcessTextractResponse(&response, testInstance.config)
		if err != nil {
			println(err.Error())
		}

	}
	t.Fail()
}
