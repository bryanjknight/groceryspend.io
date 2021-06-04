package receipts

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"math"

	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/montanaflynn/stats"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/textract"
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

// OrderedPolygonPoints ensures that the array is TL, TR, BR, BL
func OrderedPolygonPoints(pts []*textract.Point) []*textract.Point {

	// TODO: write code to ensure points are in order
	return pts
}

// Centroid calculates the centroid of a quadrilateral
func Centroid(pts []*textract.Point) *textract.Point {
	// Formula from https://en.wikipedia.org/wiki/Centroid#Of_a_polygon

	// calculate "A" first
	a := 0.0
	for i, pt := range pts {
		x := *pt.X
		y := *pt.Y
		nextX := 0.0
		nextY := 0.0
		if i == len(pts)-1 {
			nextX = *pts[0].X
			nextY = *pts[0].Y
		} else {
			nextX = *pts[i+1].X
			nextY = *pts[i+1].Y
		}

		a += (x * nextY) - (nextX * y)
	}

	// calculate cx and cy
	cx := 0.0
	cy := 0.0
	for i, pt := range pts {
		x := *pt.X
		y := *pt.Y
		nextX := 0.0
		nextY := 0.0
		if i == len(pts)-1 {
			nextX = *pts[0].X
			nextY = *pts[0].Y
		} else {
			nextX = *pts[i+1].X
			nextY = *pts[i+1].Y
		}

		cx += (x + nextX) * (x*nextY - nextX*y)
		cy += (y + nextY) * (x*nextY - nextX*y)
	}

	// finalize cx and cy
	cx = cx / (6 * a)
	cy = cy / (6 * a)

	return &textract.Point{
		X: aws.Float64(cx),
		Y: aws.Float64(cy),
	}
}

// DistanceFromTwoPoints calculates the distance between two points
func DistanceFromTwoPoints(a textract.Point, b textract.Point) float64 {
	// wow, I'm finally using pythagorean theorem
	return math.Sqrt(math.Pow((*a.X-*b.X), 2) + math.Pow((*a.Y-*b.Y), 2))
}

// PolygonFromBlocks will return an array of TL TR, BR, BL
func PolygonFromBlocks(blocks []*textract.Block) []*textract.Point {

	topLeft := textract.Point{X: aws.Float64(0.0), Y: aws.Float64(0.0)}
	bottomLeft := textract.Point{X: aws.Float64(0.0), Y: aws.Float64(1.0)}
	topRight := textract.Point{X: aws.Float64(1.0), Y: aws.Float64(0.0)}
	bottomRight := textract.Point{X: aws.Float64(1.0), Y: aws.Float64(1.0)}

	// we set all of these to impossible values so that
	// as we iterate through the blocks, they converage on the right solution
	regionTopLeft := textract.Point{X: aws.Float64(1.0), Y: aws.Float64(1.0)}
	tlDist := DistanceFromTwoPoints(topLeft, regionTopLeft)
	regionBottomLeft := textract.Point{X: aws.Float64(1.0), Y: aws.Float64(0.0)}
	blDist := DistanceFromTwoPoints(bottomLeft, regionBottomLeft)
	regionTopRight := textract.Point{X: aws.Float64(0.0), Y: aws.Float64(1.0)}
	trDist := DistanceFromTwoPoints(topRight, regionTopRight)
	regionBottomRight := textract.Point{X: aws.Float64(0.0), Y: aws.Float64(0.0)}
	brDist := DistanceFromTwoPoints(bottomRight, regionBottomRight)

	for _, block := range blocks {
		pts := OrderedPolygonPoints(block.Geometry.Polygon)

		possibleTopLeft := pts[0]
		possibleTlDist := DistanceFromTwoPoints(topLeft, *possibleTopLeft)
		possibleTopRight := pts[1]
		possibleTrDist := DistanceFromTwoPoints(topLeft, *possibleTopRight)
		possibleBottomRight := pts[2]
		possibleBrDist := DistanceFromTwoPoints(topLeft, *possibleBottomRight)
		possibleBottomLeft := pts[3]
		possibleBlDist := DistanceFromTwoPoints(topLeft, *possibleTopRight)

		if possibleTlDist < tlDist {
			regionTopLeft = *possibleTopLeft
			tlDist = possibleTlDist
		}

		if possibleTrDist < trDist {
			regionTopRight = *possibleTopRight
			trDist = possibleTrDist
		}

		if possibleBrDist < brDist {
			regionBottomRight = *possibleBottomRight
			brDist = possibleBrDist
		}

		if possibleBlDist < blDist {
			regionBottomLeft = *possibleBottomLeft
			blDist = possibleBlDist
		}

	}

	return []*textract.Point{
		&regionTopLeft,
		&regionTopRight,
		&regionBottomRight,
		&regionBottomLeft,
	}
}

// RegressionIntersectsWithSegment determines if a linear regression passes through a line segment
func RegressionIntersectsWithSegment(lr *linearRegression, a *textract.Point, b *textract.Point) bool {
	// see https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	if lr == nil || a == nil || b == nil {
		return false
	}

	x1 := *a.X
	y1 := *a.Y
	x2 := *b.X
	y2 := *b.Y

	x3 := x1
	y3 := calculateExpectedY(x1, lr)
	x4 := x2
	y4 := calculateExpectedY(x2, lr)

	d := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)

	if d == 0 {
		return false
	}

	px := ((x1*y2-y1*x2)*(x3-x4) - (x1-x2)*(x3*y4-y3*x4)) / d

	// does pX exist on the line segment ab?
	minX, _ := stats.Min([]float64{x1, x2})
	maxX, _ := stats.Max([]float64{x1, x2})
	return px >= minX && px <= maxX

}

// LinePassesThroughPolygon tests to see whether a linear regression line passes
// through a poylgon. We assume the polygon is already ordered
func LinePassesThroughPolygon(lr *linearRegression, polygon []*textract.Point) bool {

	// shoelace through the points to test each edge. If any succeed, then return true
	for i := 0; i < len(polygon); i++ {
		a := polygon[i]
		var b *textract.Point
		if i == len(polygon)-1 {
			b = polygon[0]
		} else {
			b = polygon[i+1]
		}

		// does regression line pass through edge?
		if RegressionIntersectsWithSegment(lr, a, b) {
			return true
		}
	}

	return false
}

// PolygonBetweenRegressionLines tests to see if a polygon lies between two regression lines
// we assume the polygon is already ordered
func PolygonBetweenRegressionLines(topLine *linearRegression, bottomLine *linearRegression, polygon []*textract.Point) bool {

	// get vertices of polygon
	tl := polygon[0]
	tr := polygon[1]
	br := polygon[2]
	bl := polygon[3]

	underTopLine := calculateExpectedY(*tl.X, topLine) < *tl.Y && calculateExpectedY(*tr.X, topLine) < *tr.Y
	aboveBottomLine := calculateExpectedY(*bl.X, bottomLine) > *bl.Y && calculateExpectedY(*br.X, topLine) > *br.Y

	return underTopLine && aboveBottomLine
}

func findBlocksByLinearSlope(
	possibleBlocks []*textract.Block,
	topLine *linearRegression,
	bottomLine *linearRegression,
	config *ImageReceiptParseConfig) []*textract.Block {

	retval := []*textract.Block{}

	for _, block := range possibleBlocks {

		polygon := OrderedPolygonPoints(block.Geometry.Polygon)

		intersectTopLine := LinePassesThroughPolygon(topLine, polygon)
		intersectBottomLine := LinePassesThroughPolygon(bottomLine, polygon)
		betweenBothLines := PolygonBetweenRegressionLines(topLine, bottomLine, polygon)

		// if this is a money one, just return this
		if betweenBothLines || (intersectBottomLine && intersectTopLine) {
			return []*textract.Block{block}
		} else if intersectBottomLine || intersectTopLine {
			retval = append(retval, block)
		}
	}

	return retval

}
