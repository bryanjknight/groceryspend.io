package parser

import (
	"fmt"
	"math"
	"strings"

	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/montanaflynn/stats"
	"groceryspend.io/server/utils"
)

type positionStats struct {
	mean   float64
	median float64
	stdDev float64
}

type summarySection struct {
	subTotalBlock *textract.Block
	taxBlock      *textract.Block
	totalBlock    *textract.Block
}

func (s *summarySection) String() {
	b := strings.Builder{}
	if s.subTotalBlock != nil {
		b.WriteString(fmt.Sprintf("Sub Total: %s\n", *s.subTotalBlock.Text))
	}
	if s.taxBlock != nil {
		b.WriteString(fmt.Sprintf("Tax: %s\n", *s.taxBlock.Text))
	}
	if s.totalBlock != nil {
		b.WriteString(fmt.Sprintf("Total: %s\n", *s.totalBlock.Text))
	}
}

func polygonToXpos(points []*textract.Point) []float64 {
	if points == nil {
		return nil
	}

	retval := []float64{}
	for _, pt := range points {
		retval = append(retval, *pt.X)
	}

	return retval
}

func polygonToYpos(points []*textract.Point) []float64 {
	if points == nil {
		return nil
	}

	retval := []float64{}
	for _, pt := range points {
		retval = append(retval, *pt.Y)
	}

	return retval
}

func (s *summarySection) minX() (float64, error) {
	if s.subTotalBlock == nil && s.taxBlock == nil && s.totalBlock == nil {
		return 0.0, fmt.Errorf("subtotal, tax, and total are missing")
	}

	xPos := []float64{}
	if s.subTotalBlock != nil {
		xPos = append(xPos, polygonToXpos(s.subTotalBlock.Geometry.Polygon)...)
	}
	if s.taxBlock != nil {
		xPos = append(xPos, polygonToXpos(s.taxBlock.Geometry.Polygon)...)
	}

	if s.totalBlock != nil {
		xPos = append(xPos, polygonToXpos(s.totalBlock.Geometry.Polygon)...)
	}

	return stats.Min(xPos)
}

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

func findSummarySection(resp *textract.AnalyzeDocumentOutput) (*summarySection, error) {
	// typically the format is <subtotal | tax | total> (whitespace) <total value>
	// so, we will look for the the words, then verify the current line or next line is the
	// actual value
	retval := summarySection{}

	for _, block := range resp.Blocks {
		if *block.BlockType == "LINE" && subtotalRegex.MatchString(*block.Text) {
			subTotalBlock, err := findNextLine(resp, block)
			if err != nil {
				return nil, err
			}

			retval.subTotalBlock = subTotalBlock
		} else if *block.BlockType == "LINE" && taxRegex.MatchString(*block.Text) {
			taxBlock, err := findNextLine(resp, block)
			if err != nil {
				return nil, err
			}

			retval.taxBlock = taxBlock
		} else if *block.BlockType == "LINE" && totalRegex.MatchString(*block.Text) {
			totalBlock, err := findNextLine(resp, block)
			if err != nil {
				return nil, err
			}

			retval.totalBlock = totalBlock
		}
	}

	return &retval, nil
}

// we assume the header always starts at 0,0
func findHeaderRegion(resp *textract.AnalyzeDocumentOutput) *textract.Point {

	headerBottomRight := textract.Point{
		X: new(float64), Y: new(float64),
	}

	inHeaderDetails := false

	isHeaderData := func(line string) bool {
		return dateRegex.MatchString(line) ||
			timeRegex.MatchString(line) ||
			addressRegex.MatchString(line) ||
			townCityZipRegex.MatchString(line)
	}

	for _, block := range resp.Blocks {
		if *block.BlockType != "LINE" {
			continue
		}

		if isHeaderData(*block.Text) {
			inHeaderDetails = true

			if *headerBottomRight.X < *block.Geometry.BoundingBox.Top+*block.Geometry.BoundingBox.Height {
				headerBottomRight.SetX(*block.Geometry.BoundingBox.Top + *block.Geometry.BoundingBox.Height)
			}

			if *headerBottomRight.Y < *block.Geometry.BoundingBox.Left+*block.Geometry.BoundingBox.Width {
				headerBottomRight.SetY(*block.Geometry.BoundingBox.Left + *block.Geometry.BoundingBox.Width)
			}

		} else if inHeaderDetails {
			// we have it the first line outside the header info section, so return what we have
			return &headerBottomRight
		}
	}

	// something wrong happened, we should have seen something
	return &headerBottomRight

}

func findBlocksByRegion(resp *textract.AnalyzeDocumentOutput, topLeft *textract.Point, bottomRight *textract.Point) []*textract.Block {
	retval := []*textract.Block{}

	if resp == nil || topLeft == nil || bottomRight == nil {
		return nil
	}

	for _, block := range resp.Blocks {
		if *block.Geometry.BoundingBox.Top >= *topLeft.X &&
			*block.Geometry.BoundingBox.Left >= *topLeft.Y &&
			*block.Geometry.BoundingBox.Top+*block.Geometry.BoundingBox.Height <= *bottomRight.X &&
			*block.Geometry.BoundingBox.Left+*block.Geometry.BoundingBox.Width <= *bottomRight.Y {
			retval = append(retval, block)
		}
	}

	return retval
}

func findItemFinalPrices(resp *textract.AnalyzeDocumentOutput, maxTopPos float64) ([]*textract.Block, error) {
	pass1 := []*textract.Block{}
	leftPos := []float64{}

	// first find the final prices for each line item
	// split the receipt into two columns down the middle
	for _, block := range resp.Blocks {
		if *block.BlockType == textract.BlockTypeLine &&
			priceRegex.MatchString(*block.Text) &&
			*block.Geometry.BoundingBox.Top+*block.Geometry.BoundingBox.Height < maxTopPos {
			pass1 = append(pass1, block)
			leftPos = append(leftPos, *block.Geometry.BoundingBox.Left)
		}
	}

	// calc stats on prices
	mean, err := stats.Mean(leftPos)
	if err != nil {
		return nil, err
	}
	// stdDev, err := stats.StandardDeviation(leftPos)
	// if err != nil {
	// 	return nil, err
	// }

	// println(fmt.Sprintf("mean %.5f, std dev: %.5f", mean, stdDev))

	retval := []*textract.Block{}

	for _, block := range pass1 {
		// TODO: different receipts may require different approaches
		// Possible scenarios are based on how the price is justified (usually)
		// right justified), whether it has a suffix code (e.g. F, *, or W)
		// (usually yes), and whether unit price is included (also usually yes)

		if *block.Geometry.BoundingBox.Left > mean {
			retval = append(retval, block)
		}
	}

	return retval, nil
}

type linearRegression struct {
	slope        float64
	intersection float64
}

func calculateSlopes(polygon []*textract.Point) (*linearRegression, *linearRegression, error) {

	if len(polygon) != 4 {
		println(fmt.Sprintf("Only support 4 points, got %v", len(polygon)))
		return nil, nil, fmt.Errorf("Only support 4 points, got %v", len(polygon))
	}

	// FIXME: we assume the order of points, so add logic to verify this is accurate

	topLeft := polygon[0]
	topRight := polygon[1]
	bottomRight := polygon[2]
	bottomLeft := polygon[3]

	topLineSlope := (*topRight.Y - *topLeft.Y) / (*topRight.X - *topLeft.X)
	topLineIntersect := *topLeft.Y - *topLeft.X*topLineSlope

	bottomLineSlope := (*bottomRight.Y - *bottomLeft.Y) / (*bottomRight.X - *bottomLeft.X)
	bottomLineIntersect := *bottomLeft.Y - *bottomLeft.X*bottomLineSlope

	return &linearRegression{slope: topLineSlope, intersection: topLineIntersect}, &linearRegression{slope: bottomLineSlope, intersection: bottomLineIntersect}, nil

}

func findBlocksByLinearSlope(
	resp *textract.AnalyzeDocumentOutput,
	topLine *linearRegression,
	bottomLine *linearRegression,
	desiredMaxYPos float64,
	tolerance float64) []*textract.Block {

	retval := []*textract.Block{}

	// TODO: we don't need to go through all blocks, just the ones within the item range
	for _, block := range resp.Blocks {
		if *block.BlockType != "LINE" {
			continue
		}

		// get all y positions
		yPos := polygonToYpos(block.Geometry.Polygon)

		if maxYPos, _ := stats.Max(yPos); maxYPos > desiredMaxYPos {
			continue
		}

		// FIXME: we assume the order of points, so add logic to verify this is accurate
		polygon := block.Geometry.Polygon
		topLeft := polygon[0]
		// topRight := polygon[1]
		// bottomRight := polygon[2]
		bottomLeft := polygon[3]

		desiredTopLeftY := *topLeft.X*topLine.slope + topLine.intersection
		percentOffTopLeft := (desiredTopLeftY - *topLeft.Y) / desiredTopLeftY

		desiredBottomLeftY := *bottomLeft.X*bottomLine.slope + bottomLine.intersection
		percentOffBottomLeft := (desiredBottomLeftY - *bottomLeft.Y) / desiredBottomLeftY

		if math.Abs(percentOffBottomLeft) < tolerance && math.Abs(percentOffTopLeft) < tolerance {
			retval = append(retval, block)
		}

	}

	return retval

}

// ImageReceiptParseConfig - options to configure how the textract response is processed
type ImageReceiptParseConfig struct {
	maxItemDescYPos float64
	tolerance       float64
}

// ProcessTextractResponse - find a better name
func ProcessTextractResponse(resp *textract.AnalyzeDocumentOutput, config *ImageReceiptParseConfig) error {

	itemBlocks := []*textract.Block{}

	// TODO: we do multiple O(N) operations against the array of blocks
	// an optimization could be to better structure the data for easier
	// lookup; however, perf is not a driver at this point given the size of the
	// receipts are small and take very little time (<500ms)

	// we can use a tree to create an quick way to look up where the item is
	// based on the block's top location. For example:
	// t[0.0,1.0] = all blocks (where 1000 is the max pixel location)
	// t[.5,1.0] == all blocks on the bottom half
	// this could be combined with a similar tree but using the left/right position
	// then doing a union on the resulting two sets to see what actually fits within
	// those coordinates

	// find the header section
	headerBottomRight := findHeaderRegion(resp)

	// find the total/subtotal/tax sections. This will denote where
	// the items finish up (so we don't count them towards the items)

	summary, err := findSummarySection(resp)
	if err != nil {
		return err
	}

	// verify I have everything I need
	if summary.subTotalBlock == nil && summary.taxBlock == nil && summary.totalBlock == nil {
		return fmt.Errorf("missing subtotal, tax, and total")
	}

	summaryTopXPos, err := summary.minX()
	if err != nil {
		return err
	}

	// TODO: should we pass the header bottom x pos?
	finalPriceBlocks, _ := findItemFinalPrices(resp, summaryTopXPos)

	// now run through the blocks again, this time looking for potential items
	// within the bounds of the final prices
	for _, block := range resp.Blocks {
		if *block.BlockType == textract.BlockTypeLine &&

			// it's not a department name
			!departmentNamesRegex.MatchString(*block.Text) &&

			// the top edge of the box doesn't go past the header section
			utils.IsWithinTolerance(*block.Geometry.BoundingBox.Top, *headerBottomRight.X, config.tolerance) &&

			// the right edge of the box doesn't go past hte max item description position
			utils.IsWithinTolerance(*block.Geometry.BoundingBox.Left+*block.Geometry.BoundingBox.Width,
				config.maxItemDescYPos,
				config.tolerance) &&

			// the bottom edge of the box doesn't go past the summary section
			utils.IsWithinTolerance(*block.Geometry.BoundingBox.Top+*block.Geometry.BoundingBox.Height,
				summaryTopXPos, config.tolerance) {
			itemBlocks = append(itemBlocks, block)
		}
	}

	println(fmt.Sprintf("# of final prices: %v", len(finalPriceBlocks)))
	for _, p := range finalPriceBlocks {

		topLine, bottomLine, err := calculateSlopes(p.Geometry.Polygon)
		if err != nil {
			return err
		}

		// find blocks that are within some tolerance
		println(*p.Text)

		possibleItems := findBlocksByLinearSlope(resp, topLine, bottomLine, config.maxItemDescYPos, config.tolerance)
		if len(possibleItems) > 0 {
			b := []string{}

			for _, possibleItem := range possibleItems {
				b = append(b, *possibleItem.Text)
			}

			println(fmt.Sprintf("\tPossible Items: %s", strings.Join(b, " / ")))
		} else {
			println("\tNo items found")
		}

	}
	println(fmt.Sprintf("# of item lines: %v", len(itemBlocks)))

	return nil

}