package parser

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/textract"
)

func TestTextractResponse(t *testing.T) {
	var response textract.AnalyzeDocumentOutput
	filename := filepath.Join(getTestDataDir(), "marketbasket", "receipt1-apiResponse.json")
	fileText := readFileAsString(filename)
	reader := strings.NewReader(fileText)
	err := json.NewDecoder(reader).Decode(&response)
	if err != nil {
		println(err.Error())
	}

	config := ImageReceiptParseConfig{
		maxItemDescYPos: 0.7,
		tolerance:       0.01,
	}

	err = ProcessTextractResponse(&response, &config)
	if err != nil {
		println(err.Error())
	}
	t.Fail()
}
