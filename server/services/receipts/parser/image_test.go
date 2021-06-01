package parser

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/textract"
)

func TestTextractResponse(t *testing.T) {

	type test struct {
		file              string
		expectedOrderDate time.Time
		expectedTotal     float32
	}

	// FIXME: assuming EST
	loc, _ := time.LoadLocation("America/New_York")

	mustParseTime := func(f string, s string) time.Time {
		t, err := time.ParseInLocation(f, s, loc)
		if err != nil {
			panic(err)
		}
		return t
	}

	tests := []test{
		{
			file:              filepath.Join(getTestDataDir(), "marketbasket", "receipt1-apiResponse.json"),
			expectedTotal:     34.05,
			expectedOrderDate: mustParseTime("01/02/2006", "04/03/2021"),
		},
		{
			file:              filepath.Join(getTestDataDir(), "hannaford", "receipt1-apiResponse.json"),
			expectedTotal:     29.92,
			expectedOrderDate: mustParseTime("01/02/2006", "04/06/2021"),
		},
		{
			file:              filepath.Join(getTestDataDir(), "wegmans", "receipt1-apiResponse.json"),
			expectedTotal:     64.01,
			expectedOrderDate: mustParseTime("01/02/2006", "05/16/2021"),
		},
		{
			file:              filepath.Join(getTestDataDir(), "wegmans", "receipt2-apiResponse.json"),
			expectedTotal:     55.51,
			expectedOrderDate: mustParseTime("01/02/2006", "05/04/2021"),
		},
	}

	for _, testInstance := range tests {
		t.Run(testInstance.file, func(t *testing.T) {
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
			} else if !receiptDetail.OrderTimestamp.Equal(testInstance.expectedOrderDate) {
				t.Errorf("timestamps didn't match: expected %v, got %v", testInstance.expectedOrderDate, receiptDetail.OrderTimestamp)
			}

		})
	}
}
