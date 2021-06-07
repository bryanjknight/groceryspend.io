package receipts

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/textract"
	"groceryspend.io/server/utils"
)

func createMockBlock(text string, confidence float64, xyPts []float64) *textract.Block {
	polygon := []*textract.Point{}

	for i := 0; i < len(xyPts); i += 2 {
		pt := textract.Point{
			X: aws.Float64(xyPts[i]),
			Y: aws.Float64(xyPts[i+1]),
		}
		polygon = append(polygon, &pt)
	}
	return &textract.Block{
		Geometry: &textract.Geometry{
			Polygon: polygon,
		},
		BlockType:  aws.String(textract.BlockTypeLine),
		Text:       aws.String(text),
		Confidence: aws.Float64(confidence),
	}
}

func TestIntersect(t *testing.T) {

	testBlock := &textract.Block{
		Geometry: &textract.Geometry{
			BoundingBox: &textract.BoundingBox{},
			Polygon:     []*textract.Point{{X: aws.Float64(0.6366869807243347), Y: aws.Float64(0.16261744499206543)}, {X: aws.Float64(0.679876446723938), Y: aws.Float64(0.16301347315311432)}, {X: aws.Float64(0.6794444918632507), Y: aws.Float64(0.17910079658031464)}, {X: aws.Float64(0.636359691619873), Y: aws.Float64(0.17870093882083893)}},
		},
		BlockType:  aws.String(textract.BlockTypeLine),
		Text:       aws.String("9.99"),
		Confidence: aws.Float64(99.999),
	}

	headerLr := &linearRegression{
		slope:        0.008918718940121978,
		intersection: 0.12054234549731022,
	}

	summaryLr := &linearRegression{
		slope:        0.01054189482192841,
		intersection: 0.35500219741221267,
	}

	config := &ImageReceiptParseConfig{
		ocrConfidence:                 90.0,
		regressionTolerance:           0.0,
		blocksOnHeaderLineAreHeader:   true,
		blocksOnSummaryLineAreSummary: true,
		minArea:                       .5,
	}

	belowHeader := belowLinearRegressionLine(headerLr, config, !config.blocksOnHeaderLineAreHeader)(testBlock)
	aboveSummary := aboveLinearRegressionLine(summaryLr, config, !config.blocksOnSummaryLineAreSummary)(testBlock)

	if !belowHeader || !aboveSummary {
		t.Errorf("Block should be between lines but was calculated to not be: %t, %t", belowHeader, aboveSummary)
	}
}

func TestLineItemMatch(t *testing.T) {

	type test struct {
		priceBlock         *textract.Block
		possibleItemBlocks []*textract.Block
		prevLine           *linearRegression
	}

	testCases := []test{
		{
			priceBlock: createMockBlock("1.29 F", 99.99, []float64{0.635972797870636,
				0.40718337893486023,
				0.6904643774032593,
				0.40655288100242615,
				0.6904357671737671,
				0.41863521933555603,
				0.6359216570854187,
				0.41925734281539917}),
			possibleItemBlocks: []*textract.Block{
				createMockBlock("1# PEELED BABY CARRT",
					99.50421142578125, []float64{
						0.3338685929775238,
						0.4069613516330719,
						0.5099449157714844,
						0.4049155116081238,
						0.5098243951797485,
						0.4189956784248352,
						0.3336635231971741,
						0.42100998759269714}),
			},
		},
	}

	config := &ImageReceiptParseConfig{
		regressionTolerance:           0.0,
		ocrConfidence:                 90.0,
		blocksOnHeaderLineAreHeader:   false,
		blocksOnSummaryLineAreSummary: false,
		minArea:                       .5,
	}

	for tcIdx, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tcIdx), func(t1 *testing.T) {
			_, bottomLine, err := calculateSlopes(tc.priceBlock.Geometry.Polygon)
			if err != nil {
				t1.Errorf(err.Error())
				return
			}

			possibleItems := bestEffortLineItemBlocks(
				bottomLine, config, tc.possibleItemBlocks,
			)

			if len(possibleItems) != 1 {
				t1.Errorf("expected 1, got %v", len(possibleItems))
			}
		})

	}
}

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

	confidence := 80.0

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
		// {
		// 	file:              filepath.Join(getTestDataDir(), "bjs", "receipt1-apiResponse.json"),
		// 	expectedTotal:     282.43,
		// 	expectedOrderDate: mustParseTime("01/02/2006", "05/19/2021"),
		// },
	}

	for _, testInstance := range tests {
		t.Run(testInstance.file, func(t *testing.T) {
			var response textract.DetectDocumentTextOutput
			fileText := utils.ReadFileAsString(testInstance.file)
			reader := strings.NewReader(fileText)
			err := json.NewDecoder(reader).Decode(&response)
			if err != nil {
				println(err.Error())
			}

			receiptDetail, err := ParseImageReceipt(&response, testInstance.expectedTotal, confidence)
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
