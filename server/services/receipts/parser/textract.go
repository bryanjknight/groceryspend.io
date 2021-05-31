package parser

import (
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/montanaflynn/stats"
	"groceryspend.io/server/utils"
)

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
