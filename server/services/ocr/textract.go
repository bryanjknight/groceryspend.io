package ocr

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
	"groceryspend.io/server/utils"
)

// TextractService encapsulates the connection details to AWS
type TextractService struct {
	Session      *session.Session
	S3BucketName string
}

// NewTextractService creates a new Textract Service
func NewTextractService(session *session.Session) *TextractService {

	return &TextractService{
		Session:      session,
		S3BucketName: utils.GetOsValue("OCR_AWS_S3_BUCKET_NAME"),
	}

}

// DetectTextInImage informs Textract to load the S3 file and create an Image struct of the text blocks
func (t *TextractService) DetectTextInImage(filePath string) (*Image, error) {

	// Create a Textract client from just a session.
	svc := textract.New(t.Session)

	// convert repsonse to a canonical model
	resp, err := svc.DetectDocumentText(&textract.DetectDocumentTextInput{
		Document: &textract.Document{
			S3Object: &textract.S3Object{
				Bucket: &t.S3BucketName,
				Name:   &filePath,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return TextractResponseToImage(resp)

}

func textractPointToPoint(pt *textract.Point) *Point {
	return &Point{
		X: *pt.X,
		Y: *pt.Y,
	}
}

// OrderedPolygonPoints ensures that the array is TL, TR, BR, BL
func orderedPolygonPoints(pts []*textract.Point) []*textract.Point {

	// TODO: write code to ensure points are in order
	return pts
}

// TextractResponseToImage converts a textract response to an ocr image struct
func TextractResponseToImage(resp *textract.DetectDocumentTextOutput) (*Image, error) {

	newBlocks := []*Block{}
	newBlockIDs := []string{}

	for _, block := range resp.Blocks {

		// if it's not a line, skip it
		if *block.BlockType != textract.BlockTypeLine {
			continue
		}

		if block.Geometry == nil || block.Geometry.Polygon == nil {
			println(fmt.Sprintf("Failed to find polygon for block ID: %s", *block.Id))
			continue
		}

		if len(block.Geometry.Polygon) != 4 {
			println(fmt.Sprintf("expected quadrilateral, got %v side polygon", len(block.Geometry.Polygon)))
			continue
		}

		polygon := orderedPolygonPoints(block.Geometry.Polygon)

		tmpBlock := Block{
			ID:          *block.Id,
			TopLeft:     textractPointToPoint(polygon[0]),
			TopRight:    textractPointToPoint(polygon[1]),
			BottomRight: textractPointToPoint(polygon[2]),
			BottomLeft:  textractPointToPoint(polygon[3]),
			Text:        *block.Text,
			Confidence:  *block.Confidence,
		}

		newBlocks = append(newBlocks, &tmpBlock)

		newBlockIDs = append(newBlockIDs, *block.Id)
	}

	return &Image{
		Blocks:           newBlocks,
		BlockIDs:         newBlockIDs,
		OriginalResponse: *resp,
	}, nil

}
