package receipts

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/montanaflynn/stats"
	"groceryspend.io/server/utils"
)

// UploadContentToS3 will upload the request content to S3
func UploadContentToS3(session *session.Session, request ParseReceiptRequest) (string, error) {
	// mock a response if we're running locally
	if utils.GetOsValueAsBoolean("RECEIPTS_MOCK_AWS_RESPONSE") {
		return utils.GetOsValue("RECEIPTS_MOCK_AWS_RESPONSE_FILE"), nil
	} else if session == nil {
		return "", fmt.Errorf("no session provided")
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(session)

	// get the data header info (e.g. data:image/jpeg;base64,)
	// TODO: do something with this information? perhaps not send it at all
	base64key := "base64,"
	base64idx := strings.Index(request.Data, base64key)
	base64data := strings.TrimSpace(request.Data[base64idx+len(base64key):])

	byteArr, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(byteArr)

	// FIXME: assuming jpg
	s3key := fmt.Sprintf("images/%s/image.jpg", request.ID.String())

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(utils.GetOsValue("RECEIPTS_AWS_S3_BUCKET_NAME")),
		Key:    aws.String(s3key),
		Body:   reader,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	if result == nil {
		return "", fmt.Errorf("s3 upload result is null")
	}

	return s3key, nil
}

// DetectDocumentText will ask Textract to analyze the document
func DetectDocumentText(session *session.Session, s3Key string) (*textract.DetectDocumentTextOutput, error) {

	// mock a response if we're running locally
	if utils.GetOsValueAsBoolean("RECEIPTS_MOCK_AWS_RESPONSE") {
		return mockDectedDocumentText()
	} else if session == nil {
		return nil, fmt.Errorf("no session provided")
	}

	// Create a Textract client from just a session.
	svc := textract.New(session)
	bucket := utils.GetOsValue("RECEIPTS_AWS_S3_BUCKET_NAME")

	return svc.DetectDocumentText(&textract.DetectDocumentTextInput{
		Document: &textract.Document{
			S3Object: &textract.S3Object{
				Bucket: &bucket,
				Name:   &s3Key,
			},
		},
	})

}

func mockDectedDocumentText() (*textract.DetectDocumentTextOutput, error) {
	var resp textract.DetectDocumentTextOutput
	mockRespText := utils.ReadFileAsString(utils.GetOsValue("RECEIPTS_MOCK_AWS_RESPONSE_FILE"))

	reader := strings.NewReader(mockRespText)
	err := json.NewDecoder(reader).Decode(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

func findBlocksByLinearSlope(
	blocks []*textract.Block,
	topLine *linearRegression,
	bottomLine *linearRegression,
	config *ImageReceiptParseConfig) []*textract.Block {

	retval := []*textract.Block{}

	for _, block := range blocks {

		// /remove an
		xPos := polygonToXpos(block.Geometry.Polygon)
		if maxXPos, _ := stats.Max(xPos); !utils.IsLessThanWithinTolerance(
			config.maxItemDescXPos, maxXPos, config.tolerance) {
			continue
		}

		// FIXME: we assume the order of points, so add logic to verify this is accurate
		polygon := block.Geometry.Polygon
		topLeft := polygon[0]
		// topRight := polygon[1]
		// bottomRight := polygon[2]
		bottomLeft := polygon[3]

		desiredTopLeftY := *topLeft.X*topLine.slope + topLine.intersection
		desiredBottomLeftY := *bottomLeft.X*bottomLine.slope + bottomLine.intersection

		if utils.IsWithinTolerance(desiredTopLeftY, *topLeft.Y, config.tolerance) &&
			utils.IsWithinTolerance(desiredBottomLeftY, *bottomLeft.Y, config.tolerance) {
			retval = append(retval, block)
		}

	}

	return retval

}
