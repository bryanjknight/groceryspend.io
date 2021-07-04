package receipts

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/montanaflynn/stats"
	"groceryspend.io/server/services/ocr"
	"groceryspend.io/server/utils"
)

// NewReceiptImage creates a receipt image instance based on the textract response. Returns ReceiptImage, even if
// an error occurs for debugging purposes
func NewReceiptImage(resp *ocr.Image, expectedTotal float32, config *ImageReceiptParseConfig) (*ReceiptImage, error) {

	retval := ReceiptImage{}

	// find the header section
	headerRegion, headerLine := findHeaderRegion(resp, config)
	retval.HeaderSection = headerRegion

	// find the summary section
	summaryRegion, summaryLine := findSummaryRegion(resp, config)
	retval.SummarySection = summaryRegion

	// TODO: sanity checks: do any of the regions overlap?

	// find the price section
	// cheat: find all blocks that are a price between the header and summary
	//				then only include the ones that are right most
	// find the line item desc section (basically what's leftover should be
	// the line item descriptions)
	lineItemRegion, priceRegion := findLineItemAndPriceRegions(resp, config, headerLine, summaryLine)
	retval.LineItemSection = lineItemRegion
	retval.PriceSection = priceRegion

	// sanity check the data, specifically the values in the price section
	// if they're wrong, then we'll have cascade issues further downstream that we cannot
	// resolve at that stage
	subtotalCents := 0
	for _, priceBlock := range priceRegion.blocks {
		price, err := ParsePrice(priceBlock.Text)
		if err != nil {
			return &retval, fmt.Errorf("price region has a non-price value in it: %s", priceBlock.Text)
		}
		// FIXME: NEED TO CHECK IF ITS NEGATIVE
		subtotalCents += int(float64(price) * 100)
	}

	expectedTotalCents := int(expectedTotal * 100)
	if subtotalCents > expectedTotalCents {
		return &retval, fmt.Errorf("subtotal (%v cents) is greater than expected total (%v cents)", subtotalCents, expectedTotalCents)
	}

	// https://blog.taxjar.com/cities-and-states-have-the-highest-sales-tax-rates/
	maxSalesTax := 0.115
	if !utils.IsLessThanWithinTolerance(float64(expectedTotalCents), float64(subtotalCents), maxSalesTax) {
		return &retval, fmt.Errorf("subtotal (%v cents) is too small to be correct based on expected value (%v cents)", subtotalCents, expectedTotalCents)
	}

	return &retval, nil
}

func findHeaderRegion(
	resp *ocr.Image,
	config *ImageReceiptParseConfig) (*ReceiptImageSection, *ocr.LinearRegression) {

	retval := ReceiptImageSection{}
	headerDetailBlocks := []*ocr.Block{}

	inHeaderDetails := false

	isHeaderData := func(line string) bool {
		return dateRegex.MatchString(line) ||
			timeRegex.MatchString(line) ||
			addressRegex.MatchString(line) ||
			townCityZipRegex.MatchString(line) ||
			cashierRegex.MatchString(line) ||
			storeRegex.MatchString(line) ||
			phoneNumberRegex.MatchString(line)
	}
	defaultFilter := defaultBlockFilter(config)

	for _, block := range resp.Blocks {
		if !defaultFilter(block) {
			continue
		}

		if isHeaderData(block.Text) {
			inHeaderDetails = true

			headerDetailBlocks = append(headerDetailBlocks, block)

		} else if inHeaderDetails {
			// we have it the first line outside the header info section, so return what we have
			break
		}
	}

	// determine the polygon of the header details
	headerDetailsPolygon := ocr.PolygonFromBlocks(headerDetailBlocks)

	// determine the bottom line for the header
	_, bottomLineRegression, _ := calculateSlopesForPoints(headerDetailsPolygon)

	println(fmt.Sprintf("Header line %v slope, %v intercept", bottomLineRegression.Slope, bottomLineRegression.Intercept))

	// find all blocks above the bottom line regression (include items that cross it)
	filter := func(block *ocr.Block) bool {
		return aboveLinearRegressionLine(bottomLineRegression, config, config.blocksOnHeaderLineAreHeader)(block)
	}
	blocks := filterBlocks(resp.Blocks, filter)
	retval.blocks = blocks
	retval.polygon = ocr.PolygonFromBlocks(blocks)
	return &retval, bottomLineRegression
}

func findSummaryRegion(resp *ocr.Image, config *ImageReceiptParseConfig) (*ReceiptImageSection, *ocr.LinearRegression) {

	summaryDetailBlocks := []*ocr.Block{}

	inSummaryDetails := false

	isSummaryData := func(line string) bool {
		return subtotalRegex.MatchString(line) ||
			taxRegex.MatchString(line) ||
			totalRegex.MatchString(line)
	}
	defaultFilter := defaultBlockFilter(config)

	for _, block := range resp.Blocks {
		if !defaultFilter(block) {
			continue
		}

		if isSummaryData(block.Text) {
			inSummaryDetails = true

			summaryDetailBlocks = append(summaryDetailBlocks, block)

		} else if inSummaryDetails {
			// we have it the first line outside the header info section, so return what we have
			break
		}
	}

	// determine the polygon of the header details
	summaryDetailPolygon := ocr.PolygonFromBlocks(summaryDetailBlocks)

	// determine the bottom line for the header
	topLineRegression, _, _ := calculateSlopesForPoints(summaryDetailPolygon)

	println(fmt.Sprintf("Summary line %v slope, %v intercept", topLineRegression.Slope, topLineRegression.Intercept))

	// find all blocks below the top line regression (include items that cross it)
	filter := func(block *ocr.Block) bool {
		return belowLinearRegressionLine(topLineRegression, config, config.blocksOnSummaryLineAreSummary)(block)
	}
	blocks := filterBlocks(resp.Blocks, filter)

	return &ReceiptImageSection{
		blocks:  blocks,
		polygon: ocr.PolygonFromBlocks(blocks),
	}, topLineRegression
}

func findLineItemAndPriceRegions(
	resp *ocr.Image,
	config *ImageReceiptParseConfig,
	headerLine *ocr.LinearRegression,
	summaryLine *ocr.LinearRegression,
) (*ReceiptImageSection, *ReceiptImageSection) {

	// line items and prices should be between the header and summary
	filter := func(block *ocr.Block) bool {
		return aboveLinearRegressionLine(summaryLine, config, !config.blocksOnSummaryLineAreSummary)(block) &&
			belowLinearRegressionLine(headerLine, config, !config.blocksOnHeaderLineAreHeader)(block)
	}

	lineItemAndPriceBlocks := filterBlocks(resp.Blocks, filter)

	// now separate the two via regex. We will undoubtly find prices that belong in
	// the description section, but we'll deal with that later
	possiblePriceBlocks := []*ocr.Block{}
	lineItemBlocks := []*ocr.Block{}

	for _, block := range lineItemAndPriceBlocks {
		if priceRegex.MatchString(block.Text) {
			possiblePriceBlocks = append(possiblePriceBlocks, block)
		} else if !departmentNamesRegex.MatchString(block.Text) {
			lineItemBlocks = append(lineItemBlocks, block)
		}
	}

	// now we need to figure out what prices are the final prices and which are
	// unit prices. To do this, we'll create another regression line, then associate
	// anything to the right of the line (within some tolerance) as being a final price
	// we'll use the centroid of the block's polygon as the input for the lineaer regression
	centroids := []*ocr.Point{}
	coords := []stats.Coordinate{}
	blockIds := []string{}
	blockIDToBlock := make(map[string]*ocr.Block)
	for _, possiblePriceBlock := range possiblePriceBlocks {
		blockIDToBlock[possiblePriceBlock.ID] = possiblePriceBlock
		centroid := ocr.Centroid(possiblePriceBlock)
		centroids = append(centroids, centroid)
		blockIds = append(blockIds, possiblePriceBlock.ID)
		// IMPORTANT: we switch  Y and X because we want the X coord to be
		//						the output and Y as the input
		coordinate := stats.Coordinate{X: centroid.Y, Y: centroid.X}
		coords = append(coords, coordinate)
	}

	centroidLinerRegression, err := stats.LinReg(coords)
	if err != nil {
		return nil, nil
	}

	// now going through the centroid regressions, and anything to the left
	// gets put into the line item description bucket
	priceBlocks := []*ocr.Block{}
	for itr, regressionXValue := range centroidLinerRegression {
		centroid := centroids[itr]
		blockID := blockIds[itr]
		block := blockIDToBlock[blockID]
		// REMEMBER: the regression output is the x axis
		if utils.IsGreaterThanWithinTolerance(regressionXValue.Y, centroid.X, config.regressionTolerance) {
			priceBlocks = append(priceBlocks, block)
		} else {
			lineItemBlocks = append(lineItemBlocks, block)
		}
	}

	lineItemBlocks = ocr.SortBlocksByLogicalOrder(resp, lineItemBlocks)

	lineItemRegion := &ReceiptImageSection{
		blocks:  lineItemBlocks,
		polygon: ocr.PolygonFromBlocks(lineItemBlocks),
	}

	priceRegion := &ReceiptImageSection{
		blocks:  priceBlocks,
		polygon: ocr.PolygonFromBlocks(priceBlocks),
	}

	return lineItemRegion, priceRegion
}

// renderReceiptImageLines renders the various lines onto an image for better debugging
func renderReceiptImageLines(ri *ReceiptImage, originalImageFilePath string, outputFileName string) error {

	f, err := os.Open(originalImageFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	srcImg, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	imgBounds := srcImg.Bounds()
	dstImg := image.NewRGBA(imgBounds)

	// copy src to dst
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.ZP, draw.Src)

	// red
	headerColor := color.RGBA{0xff, 0x00, 0x00, 0x3F}

	// blue
	summaryColor := color.RGBA{0x00, 0x00, 0xff, 0x3F}

	// green
	lineitemColor := color.RGBA{0x00, 0xff, 0x00, 0x3F}

	// cyan
	priceColor := color.RGBA{0x00, 0xff, 0xff, 0x3F}

	type regionConfig struct {
		region *ReceiptImageSection
		color  color.RGBA
	}

	work := []regionConfig{
		{ri.HeaderSection, headerColor},
		{ri.SummarySection, summaryColor},
		{ri.LineItemSection, lineitemColor},
		{ri.PriceSection, priceColor},
	}

	for _, worktItem := range work {

		// convert from 0 to 1 to 0 to N
		tlX := int(worktItem.region.polygon[0].X * float64(imgBounds.Max.X))
		tlY := int(worktItem.region.polygon[0].Y * float64(imgBounds.Max.Y))
		brX := int(worktItem.region.polygon[2].X * float64(imgBounds.Max.X))
		brY := int(worktItem.region.polygon[2].Y * float64(imgBounds.Max.Y))

		println(fmt.Sprintf("rect: (%v,%v) (%v,%v)", tlX, tlY, brX, brY))
		rectToDraw := image.Rect(
			tlX,
			tlY,
			brX,
			brY,
		)

		// draw the rectangle
		draw.Draw(dstImg, rectToDraw, &image.Uniform{worktItem.color}, image.ZP, draw.Over)

	}

	myfile, err := os.Create(fmt.Sprintf("%s.jpg", outputFileName))
	if err != nil {
		return err
	}
	err = jpeg.Encode(myfile, dstImg, nil)
	if err != nil {
		return err
	}

	return nil
}
