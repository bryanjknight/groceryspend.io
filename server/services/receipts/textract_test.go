package receipts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/montanaflynn/stats"
)

func f(x float64, y float64) *textract.Point {
	return &textract.Point{
		X: aws.Float64(x),
		Y: aws.Float64(y),
	}
}

func makePolygon(xyPts []float64) []*textract.Point {

	retval := []*textract.Point{}
	for i := 0; i < len(xyPts); i += 2 {
		retval = append(retval, f(xyPts[i], xyPts[i+1]))
	}

	return retval
}

func polygonsAreEqual(a []*textract.Point, b []*textract.Point) (bool, error) {
	if len(a) != len(b) {
		return false, fmt.Errorf("lengths not the same, got %v and %v", len(a), len(b))
	}

	for i := 0; i < len(a); i++ {
		x1, _ := stats.Round(*a[i].X, 5)
		y1, _ := stats.Round(*a[i].Y, 5)
		x2, _ := stats.Round(*b[i].X, 5)
		y2, _ := stats.Round(*b[i].Y, 5)
		if x1 != x2 {
			return false, fmt.Errorf("i: %v, Ax: %v, Bx: %v", i, *a[i].X, *b[i].X)
		}
		if y1 != y2 {
			return false, fmt.Errorf("i: %v, Ay: %v, By: %v", i, *a[i].Y, *b[i].Y)
		}
	}

	return true, nil
}

func polygonToString(polygon []*textract.Point) string {
	buffer := strings.Builder{}

	for _, pt := range polygon {
		buffer.WriteString(fmt.Sprintf("(%v, %v) ", *pt.X, *pt.Y))
	}

	return buffer.String()
}

func TestPolygonFromBlocks(t *testing.T) {
	blocks := []*textract.Block{
		{
			Geometry: &textract.Geometry{
				Polygon: []*textract.Point{
					{X: aws.Float64(1), Y: aws.Float64(1)},
					{X: aws.Float64(2), Y: aws.Float64(1)},
					{X: aws.Float64(2), Y: aws.Float64(2)},
					{X: aws.Float64(1), Y: aws.Float64(2)},
				},
			},
		},
		{
			Geometry: &textract.Geometry{
				Polygon: []*textract.Point{
					{X: aws.Float64(1), Y: aws.Float64(3)},
					{X: aws.Float64(3), Y: aws.Float64(3)},
					{X: aws.Float64(4), Y: aws.Float64(4)},
					{X: aws.Float64(1), Y: aws.Float64(4)},
				},
			},
		},
	}
	expectedPolygon := []*textract.Point{
		{X: aws.Float64(1), Y: aws.Float64(1)},
		{X: aws.Float64(4), Y: aws.Float64(1)},
		{X: aws.Float64(4), Y: aws.Float64(4)},
		{X: aws.Float64(1), Y: aws.Float64(4)},
	}
	polygon := PolygonFromBlocks(blocks)
	for i := 0; i < len(polygon); i++ {
		expectedPt := expectedPolygon[i]
		actualPt := polygon[i]

		if *actualPt.X != *expectedPt.X || *actualPt.Y != *expectedPt.Y {
			t.Errorf("Expected (%v, %v) got (%v, %v) for index %v",
				*expectedPt.X, *expectedPt.Y, *actualPt.X, *actualPt.Y, i)
		}
	}

}

func TestPolygonArea(t *testing.T) {
	type test struct {
		polygon      []*textract.Point
		expectedArea float64
	}

	testCases := []test{
		{
			polygon:      makePolygon([]float64{1, 0, 1, 1, 0, 1, 0, 0}),
			expectedArea: 1,
		},
		{
			polygon:      makePolygon([]float64{-1, 1, 1, 1, 1, -1, -1, -1}),
			expectedArea: 4,
		},
	}

	for tcItr, tc := range testCases {
		t.Run(fmt.Sprintf("Test Case %v", tcItr), func(t *testing.T) {
			actualArea := PolygonArea(tc.polygon)
			if actualArea != tc.expectedArea {
				t.Errorf("Expected (%v), got (%v)", tc.expectedArea, actualArea)
			}
		})
	}

}

func TestRegressionSegmentIntersect(t *testing.T) {
	type test struct {
		lr            *linearRegression
		a             *textract.Point
		b             *textract.Point
		doesIntersect bool
	}

	testCases := []test{
		// intersects at 1,1
		{
			lr:            &linearRegression{slope: 0, intersection: 1},
			a:             f(0, 0),
			b:             f(2, 2),
			doesIntersect: true,
		},
		// intersects at 2,2
		{
			lr:            &linearRegression{slope: 0.5, intersection: 1},
			a:             f(0, 0),
			b:             f(2, 2),
			doesIntersect: true,
		},
		// lines are parallel
		{
			lr:            &linearRegression{slope: 1, intersection: 1},
			a:             f(0, 0),
			b:             f(2, 2),
			doesIntersect: false,
		},
		// lines intersect outside line segment
		{
			lr:            &linearRegression{slope: 1, intersection: 1},
			a:             f(0, 0),
			b:             f(2, -2),
			doesIntersect: false,
		},
	}

	for tcItr, tc := range testCases {
		t.Run(fmt.Sprintf("Test Case %v", tcItr), func(t1 *testing.T) {
			pt := IntersectionBetweenRegressionAndSegment(tc.lr, tc.a, tc.b)
			if PointExistsOnLine(tc.a, tc.b, pt) != tc.doesIntersect {
				t1.Errorf("Expected (%t), got (%t)", tc.doesIntersect, PointExistsOnLine(tc.a, tc.b, pt))
			}
		})
	}
}

func TestPolygonsCreatedByCrossingLine(t *testing.T) {
	type test struct {
		polygon       []*textract.Point
		lr            *linearRegression
		topPolygon    []*textract.Point
		bottomPolygon []*textract.Point
	}

	testCases := []test{
		{
			polygon:       makePolygon([]float64{.1, .1, .5, .09, .5, .4, .9, .39}),
			lr:            &linearRegression{slope: 0, intersection: .2},
			topPolygon:    makePolygon([]float64{.1, .1, .5, .09, .5, .2, .37586, .2}),
			bottomPolygon: makePolygon([]float64{.37586, .2, .5, .2, .5, .4, .9, .39}),
		},
		{
			polygon:       makePolygon([]float64{1, 1, 5, 1, 5, 5, 5, 1}),
			lr:            &linearRegression{slope: 1, intersection: 0},
			topPolygon:    makePolygon([]float64{1, 1, 5, 1, 5, 5}),
			bottomPolygon: makePolygon([]float64{1, 1, 5, 5, 5, 1}),
		},
	}

	for itr, tc := range testCases {
		t.Run(fmt.Sprintf("%v", itr), func(t1 *testing.T) {
			top, bottom, err := PolygonsCreatedByCrossingLine(tc.polygon, tc.lr)
			if err != nil {
				t1.Errorf("Failed to calculate polygons: %s", err.Error())
				return
			}

			// check top
			if ok, err := polygonsAreEqual(top, tc.topPolygon); !ok {
				t1.Errorf("top polygons did not match: %s", err.Error())
				return
			}

			//check bottom
			if ok, err := polygonsAreEqual(bottom, tc.bottomPolygon); !ok {
				t1.Errorf("bottom polygons did not match: %s", err.Error())
				return
			}
		})
	}
}
