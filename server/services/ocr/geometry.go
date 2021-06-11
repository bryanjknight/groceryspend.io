package ocr

import (
	"fmt"
	"math"

	"github.com/montanaflynn/stats"
)

// BlockToPointArray converts a block's points to an array of points consisting of the following order:
// TL, TR, BR, BL. This is to make the refactored code backwards compatible with the legacy code used
// for textract
func BlockToPointArray(block *Block) []*Point {
	return []*Point{
		block.TopLeft, block.TopRight, block.BottomRight, block.BottomLeft,
	}
}

// Centroid calculates the centroid of a quadrilateral
func Centroid(block *Block) *Point {
	// Formula from https://en.wikipedia.org/wiki/Centroid#Of_a_polygon

	pts := BlockToPointArray(block)

	// calculate "A" first
	a := 0.0
	for i, pt := range pts {
		x := pt.X
		y := pt.Y
		nextX := 0.0
		nextY := 0.0
		if i == len(pts)-1 {
			nextX = pts[0].X
			nextY = pts[0].Y
		} else {
			nextX = pts[i+1].X
			nextY = pts[i+1].Y
		}

		a += (x * nextY) - (nextX * y)
	}

	// calculate cx and cy
	cx := 0.0
	cy := 0.0
	for i, pt := range pts {
		x := pt.X
		y := pt.Y
		nextX := 0.0
		nextY := 0.0
		if i == len(pts)-1 {
			nextX = pts[0].X
			nextY = pts[0].Y
		} else {
			nextX = pts[i+1].X
			nextY = pts[i+1].Y
		}

		cx += (x + nextX) * (x*nextY - nextX*y)
		cy += (y + nextY) * (x*nextY - nextX*y)
	}

	// finalize cx and cy
	cx = cx / (6 * a)
	cy = cy / (6 * a)

	return &Point{
		X: cx,
		Y: cy,
	}
}

// DistanceFromTwoPoints calculates the distance between two points
func DistanceFromTwoPoints(a *Point, b *Point) float64 {
	// wow, I'm finally using pythagorean theorem
	return math.Sqrt(math.Pow((a.X-b.X), 2) + math.Pow((a.Y-b.Y), 2))
}

// PolygonFromBlocks will return an array of TL TR, BR, BL. it will take the min and max
// values of the x and y positions to create a region that encompases all the blocks in the array
func PolygonFromBlocks(blocks []*Block) []*Point {

	allXpos := []float64{}
	allYpos := []float64{}
	for _, block := range blocks {
		pts := BlockToPointArray(block)
		for _, point := range pts {
			allXpos = append(allXpos, point.X)
			allYpos = append(allYpos, point.Y)
		}
	}

	minX, _ := stats.Min(allXpos)
	maxX, _ := stats.Max(allXpos)
	minY, _ := stats.Min(allYpos)
	maxY, _ := stats.Max(allYpos)

	return []*Point{
		{X: minX, Y: minY},
		{X: maxX, Y: minY},
		{X: maxX, Y: maxY},
		{X: minX, Y: maxY},
	}

}

// IntersectionBetweenRegressionAndSegment determines if a linear regression passes through a line segment
func IntersectionBetweenRegressionAndSegment(lr *LinearRegression, a *Point, b *Point) *Point {
	// see https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	if lr == nil || a == nil || b == nil {
		return nil
	}

	x1 := a.X
	y1 := a.Y
	x2 := b.X
	y2 := b.Y

	x3 := x1
	y3 := CalculateExpectedY(x1, lr)
	x4 := x2
	y4 := CalculateExpectedY(x2, lr)

	d := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)

	var px, py float64

	// if the regression line is 0 sloped
	if x3 == x4 && y3 == y4 {
		px = x3
		py = y3
	} else if d == 0 {
		return nil
	} else {
		px = ((x1*y2-y1*x2)*(x3-x4) - (x1-x2)*(x3*y4-y3*x4)) / d
		py = ((x1*y2-y1*x2)*(y3-y4) - (y1-y2)*(x3*y4-y3*x4)) / d

	}

	//  check that px is between the two x values of the line
	minX, _ := stats.Min([]float64{x1, x2})
	maxX, _ := stats.Max([]float64{x1, x2})
	if px < minX || px > maxX {
		return nil
	}

	return &Point{
		X: px,
		Y: py,
	}
}

// PointExistsOnLine determines if a point exists on a line. We round to 5 decimal points
// to avoid issues with float calculations
func PointExistsOnLine(a *Point, b *Point, pt *Point) bool {

	if a == nil || b == nil || pt == nil {
		return false
	}

	x1 := a.X
	y1 := a.Y

	x2 := b.X
	y2 := b.Y

	//  check that px is between the two x values of the line
	minX, _ := stats.Min([]float64{x1, x2})
	maxX, _ := stats.Max([]float64{x1, x2})

	if pt.X < minX || pt.X > maxX {
		return false
	}

	slope := (y2 - y1) / (x2 - x1)
	intercept := y1 - x1*slope

	expectedY, _ := stats.Round(pt.X*slope+intercept, 5)
	actualY, _ := stats.Round(pt.Y, 5)

	return expectedY == actualY
}

// TriangleArea calcuates the area of a triangle given the three sides of the triangle
func TriangleArea(ab float64, bc float64, ca float64) float64 {
	// https://en.wikipedia.org/wiki/Heron%27s_formula
	s := (ab + bc + ca) / 2
	a := math.Sqrt(
		s * (s - ab) * (s - bc) * (s - ca))

	return a
}

// PolygonArea calculates the area of a polygon
func PolygonArea(polygon []*Point) float64 {

	tl := polygon[0]
	tr := polygon[1]
	tltr := DistanceFromTwoPoints(tl, tr)
	br := polygon[2]
	trbr := DistanceFromTwoPoints(tr, br)

	if len(polygon) == 3 {
		tlbr := DistanceFromTwoPoints(tl, br)
		return TriangleArea(tltr, trbr, tlbr)
	}

	bl := polygon[3]

	tlbl := DistanceFromTwoPoints(tl, bl)
	bltr := DistanceFromTwoPoints(bl, tr)
	blbr := DistanceFromTwoPoints(bl, br)

	// do it twice for the two triangles in the polygon

	triArea1 := TriangleArea(tlbl, tltr, bltr)
	triArea2 := TriangleArea(bltr, blbr, trbr)

	roundedArea, _ := stats.Round(triArea1+triArea2, 5)
	return roundedArea
}

// PercentagePolygonCoveredByLine determine the area covered by the linear regression. A flag 'isAbove'
// determines if the area is above the regression line. Min area is a value from 0 to 1
func PercentagePolygonCoveredByLine(lr *LinearRegression, polygon []*Point, isAbove bool) (float64, error) {

	// 1) determine points where line passes through polygon. We can do this by using the
	//    the formula in RegressionIntersectsWithSegment
	top, bottom, err := PolygonsCreatedByCrossingLine(polygon, lr)
	if err != nil {
		return 0.0, err
	}

	// 2) determine area of polygon
	polygonArea := PolygonArea(polygon)

	// 3) determine area of intersected area. By figuring out the triangle, we can use that or subtract
	//    from the total area to get the final area
	if isAbove && len(top) <= 4 {
		// if we want the top polygon's percentage, and we can caculate the
		// top polygon
		return PolygonArea(top) / polygonArea, nil
	} else if isAbove && len(top) > 4 {
		// if we want above, but the above poly is a complex poly, just remove the bottom part from the overall area
		return (polygonArea - PolygonArea(bottom)) / polygonArea, nil
	} else if !isAbove && len(top) <= 4 {
		// if we want below, and the top is simple, then remove the top
		return (polygonArea - PolygonArea(top)) / polygonArea, nil
	} else if !isAbove && len(top) > 5 {
		// if we want below the line, and the top is complex, then just use the bottom
		return PolygonArea(bottom) / polygonArea, nil
	}
	return 0.0, fmt.Errorf("unexpected error with top/bottom polygon area calculation")

}

// CompactPolygon removes any duplicate points
func CompactPolygon(polygon []*Point) []*Point {

	trackingSet := make(map[string]bool)

	retval := []*Point{}

	for _, pt := range polygon {
		key := fmt.Sprintf("%.5f,%.5f", pt.X, pt.Y)
		if _, ok := trackingSet[key]; !ok {
			retval = append(retval, pt)
			trackingSet[key] = true
		}
	}

	return retval

}

// PolygonsCreatedByCrossingLine returns two polygons created by a crossing regression line. The first returned value
// is the polygon created above the line, and the second is the polygon below the line. Recall that the
// perspective is a receipt, so 0,0 is the top-left corner of the image, X goes left to right, y goes from top to bottom;
// therefore, "top" is would be anything going to Y->0
func PolygonsCreatedByCrossingLine(originalPolygon []*Point, lr *LinearRegression) ([]*Point, []*Point, error) {

	if len(originalPolygon) != 4 {
		return nil, nil, fmt.Errorf("only support 4 sided polygon, got %v", len(originalPolygon))
	}

	// we will create two arrays, which will track the 8 possible points that will come out of
	// calculating the two polygons
	trackingArray := make([]*Point, 8)

	// shoelace through the points to test each edge
	// we can precompute the possible combinations given we only support a 4 sided polygon
	// Note (N) is the index in the tracking array
	//    TL(0) ---- A(1) ---- TR(2)
	//    |                      |
	//    D(7)                  B(3)
	//    |                      |
	//    BL(6) ---- C(5) ---- BR(4)
	//
	trackingArray[0] = originalPolygon[0]
	trackingArray[2] = originalPolygon[1]
	trackingArray[4] = originalPolygon[2]
	trackingArray[6] = originalPolygon[3]

	trackingArray[1] = IntersectionBetweenRegressionAndSegment(lr, trackingArray[0], trackingArray[2])
	trackingArray[3] = IntersectionBetweenRegressionAndSegment(lr, trackingArray[2], trackingArray[4])
	trackingArray[5] = IntersectionBetweenRegressionAndSegment(lr, trackingArray[4], trackingArray[6])
	trackingArray[7] = IntersectionBetweenRegressionAndSegment(lr, trackingArray[6], trackingArray[0])

	// Crosses A and B
	if trackingArray[1] != nil && trackingArray[3] != nil {
		return []*Point{
				trackingArray[1],
				trackingArray[2],
				trackingArray[3],
			}, CompactPolygon([]*Point{
				trackingArray[0],
				trackingArray[1],
				trackingArray[3],
				trackingArray[4],
				trackingArray[6],
			}),
			nil
	} else if trackingArray[1] != nil && trackingArray[3] != nil {
		// crosses A and C
		left := []*Point{
			trackingArray[0],
			trackingArray[1],
			trackingArray[5],
			trackingArray[6],
		}
		right := []*Point{
			trackingArray[1],
			trackingArray[2],
			trackingArray[4],
			trackingArray[5],
		}

		// if the slope is negative, the right will be the "top"
		// FIXME: we don't support vertical lines
		if lr.Slope < 0 {
			return right, left, nil
		}
		return left, right, nil
	} else if trackingArray[1] != nil && trackingArray[7] != nil {
		// crosses AD
		return []*Point{
				trackingArray[1],
				trackingArray[2],
				trackingArray[7],
			}, CompactPolygon([]*Point{
				trackingArray[1],
				trackingArray[2],
				trackingArray[4],
				trackingArray[6],
				trackingArray[7],
			}),
			nil
	} else if trackingArray[3] != nil && trackingArray[5] != nil {
		// crosses BC
		return CompactPolygon([]*Point{
				trackingArray[0],
				trackingArray[2],
				trackingArray[3],
				trackingArray[5],
				trackingArray[6],
			}), []*Point{
				trackingArray[3],
				trackingArray[4],
				trackingArray[5],
			},
			nil
	} else if trackingArray[3] != nil && trackingArray[7] != nil {
		// crosses BD
		return []*Point{
				trackingArray[0],
				trackingArray[2],
				trackingArray[3],
				trackingArray[7],
			}, []*Point{
				trackingArray[7],
				trackingArray[3],
				trackingArray[4],
				trackingArray[6],
			},
			nil
	} else if trackingArray[5] != nil && trackingArray[7] != nil {
		// crosses CD
		return CompactPolygon([]*Point{
				trackingArray[0],
				trackingArray[2],
				trackingArray[4],
				trackingArray[5],
				trackingArray[7],
			}), []*Point{
				trackingArray[7],
				trackingArray[5],
				trackingArray[6],
			},
			nil
	}

	// if we get here, then there was no intersection
	return nil, nil, fmt.Errorf("No intersection found")

}

// LinePassesThroughPolygon tests to see whether a linear regression line passes
// through a poylgon. We assume the polygon is already ordered
func LinePassesThroughPolygon(lr *LinearRegression, polygon []*Point) bool {

	// shoelace through the points to test each edge. If any succeed, then return true
	for i := 0; i < len(polygon); i++ {
		a := polygon[i]
		var b *Point
		if i == len(polygon)-1 {
			b = polygon[0]
		} else {
			b = polygon[i+1]
		}

		// does regression line pass through edge?
		pt := IntersectionBetweenRegressionAndSegment(lr, a, b)
		if PointExistsOnLine(a, b, pt) {
			return true
		}
	}

	return false
}

// PolygonBetweenRegressionLines tests to see if a polygon lies between two regression lines
// we assume the polygon is already ordered
func PolygonBetweenRegressionLines(topLine *LinearRegression, bottomLine *LinearRegression, polygon []*Point) bool {

	// get vertices of polygon
	tl := polygon[0]
	tr := polygon[1]
	br := polygon[2]
	bl := polygon[3]

	underTopLine := CalculateExpectedY(tl.X, topLine) < tl.Y && CalculateExpectedY(tr.X, topLine) < tr.Y
	aboveBottomLine := CalculateExpectedY(bl.X, bottomLine) > bl.Y && CalculateExpectedY(br.X, topLine) > br.Y

	return underTopLine && aboveBottomLine
}

// FindBlocksByLinearSlope findes blocks that that either are between the two regerssion lines, or intersect more than
// the min area required to be part of that regression line
func FindBlocksByLinearSlope(
	possibleBlocks []*Block,
	topLine *LinearRegression,
	bottomLine *LinearRegression,
	minArea float64) []*Block {

	retval := []*Block{}

	for _, block := range possibleBlocks {

		polygon := BlockToPointArray(block)

		meetsTopLineMin := false
		topLineArea, tlErr := PercentagePolygonCoveredByLine(topLine, polygon, false)
		if tlErr == nil && topLineArea >= minArea {
			meetsTopLineMin = true
		}
		meetsBottomLineMin := false
		bottomLineArea, blErr := PercentagePolygonCoveredByLine(bottomLine, polygon, true)
		if blErr == nil && bottomLineArea >= minArea {
			meetsBottomLineMin = true
		}
		betweenBothLines := PolygonBetweenRegressionLines(topLine, bottomLine, polygon)

		// if this is a money one, just return this
		if betweenBothLines || meetsTopLineMin || meetsBottomLineMin {
			return []*Block{block}
		} else if (tlErr == nil && topLineArea > 0) || (blErr == nil && bottomLineArea > 0) {
			retval = append(retval, block)
		}
	}

	return retval

}
