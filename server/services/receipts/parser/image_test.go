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
		t.Fail()
	}

	err = ProcessTextractResponse(&response)
	if err != nil {
		println(err.Error())
		t.Fail()
	}
	t.Fail()
}
