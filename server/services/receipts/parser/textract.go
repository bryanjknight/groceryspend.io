package parser

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/montanaflynn/stats"
	"groceryspend.io/server/utils"
)

func findPages(resp *textract.AnalyzeDocumentOutput) []*textract.Block {
	retval := []*textract.Block{}

	for _, block := range resp.Blocks {
		if *block.BlockType == "PAGE" {
			retval = append(retval, block)
		}
	}
	return retval
}

func findPageForLine(resp *textract.AnalyzeDocumentOutput, lineBlock *textract.Block) *textract.Block {
	// get the pages
	pages := findPages(resp)

	// TODO: optimize by populating the resp into a better data structure
	for _, page := range pages {
		for _, relation := range page.Relationships {
			if *relation.Type != "CHILD" {
				continue
			}
			for _, id := range relation.Ids {
				if id == lineBlock.Id {
					return page
				}
			}
		}
	}

	return nil
}

func findBlockByID(resp *textract.AnalyzeDocumentOutput, id string) (*textract.Block, error) {
	for _, block := range resp.Blocks {
		if *block.Id == id {
			return block, nil
		}
	}

	return nil, fmt.Errorf("did not find a block with ID %s", id)
}

func findBlockByPageIDandIdx(resp *textract.AnalyzeDocumentOutput, pageID string, childIdx int) (*textract.Block, error) {
	page, err := findBlockByID(resp, pageID)
	if err != nil {
		return nil, err
	}

	for _, relation := range page.Relationships {
		if *relation.Type != "CHILD" {
			continue
		}

		return findBlockByID(resp, *relation.Ids[childIdx])
	}

	return nil, fmt.Errorf("did not find block by page id %s and child idx %v", pageID, childIdx)
}

func findNextLine(resp *textract.AnalyzeDocumentOutput, lineBlock *textract.Block) (*textract.Block, error) {
	// get the pages
	pages := findPages(resp)

	// TODO: optimize by populating the resp into a better data structure
	for pageIdx, page := range pages {
		for _, relation := range page.Relationships {
			if *relation.Type != "CHILD" {
				continue
			}
			for childIdx, id := range relation.Ids {
				if *id != *lineBlock.Id {
					continue
				}

				// if there's still children left on the page
				if childIdx < len(relation.Ids)-1 {
					nextLineID := relation.Ids[childIdx+1]
					return findBlockByID(resp, *nextLineID)
				} else if pageIdx < len(pages)-1 {

					// FIXME: this will break if the 1st child on the next page is not a line
					return findBlockByPageIDandIdx(resp, *pages[pageIdx+1].Id, 0)
				}
			}
		}
	}

	return nil, fmt.Errorf("Did not find a next line for %s", *lineBlock.Id)
}

func findBlocksByRegion(resp *textract.AnalyzeDocumentOutput, topLeft *textract.Point, bottomRight *textract.Point) []*textract.Block {
	retval := []*textract.Block{}

	if resp == nil || topLeft == nil || bottomRight == nil {
		return nil
	}

	for _, block := range resp.Blocks {
		if *block.Geometry.BoundingBox.Top >= *topLeft.Y &&
			*block.Geometry.BoundingBox.Left >= *topLeft.X &&
			*block.Geometry.BoundingBox.Top+*block.Geometry.BoundingBox.Height <= *bottomRight.Y &&
			*block.Geometry.BoundingBox.Left+*block.Geometry.BoundingBox.Width <= *bottomRight.X {
			retval = append(retval, block)
		}
	}

	return retval
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
