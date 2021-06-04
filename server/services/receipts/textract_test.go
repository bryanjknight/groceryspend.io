package receipts

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/textract"
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

func TestCentroid(t *testing.T) {
	type test struct {
		polygon          []*textract.Point
		expectedCentroid *textract.Point
	}

	testCases := []test{
		{
			polygon:          makePolygon([]float64{-1, 1, 1, 1, 1, -1, -1, -1}),
			expectedCentroid: f(0, 0),
		},
		{
			polygon:          makePolygon([]float64{2, 0, 2, 2, 0, 2, 0, 0}),
			expectedCentroid: f(0.5, 0.5),
		},
	}

	for _, tc := range testCases {
		actualCentroid := Centroid(tc.polygon)
		if *actualCentroid.Y != *tc.expectedCentroid.Y ||
			*actualCentroid.X != *tc.expectedCentroid.X {
			t.Errorf("Expected (%v, %v), got (%v, %v)",
				*tc.expectedCentroid.X, *tc.expectedCentroid.Y,
				*actualCentroid.X, *actualCentroid.Y)
		}
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
		t.Run(fmt.Sprintf("Test Case %v", tcItr), func(t *testing.T) {
			intersected := RegressionIntersectsWithSegment(tc.lr, tc.a, tc.b)
			if intersected != tc.doesIntersect {
				t.Errorf("Expected (%t), got (%t)", tc.doesIntersect, intersected)
			}
		})
	}
}
