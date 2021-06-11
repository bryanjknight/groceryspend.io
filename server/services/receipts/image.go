package receipts

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
	"groceryspend.io/server/services/ocr"
	"groceryspend.io/server/utils"
)

// ReceiptImageSection is a general region of a receipt. Typically there are 4 parts:
// a header at the top, the line items, and the summary at the bottom. The line items are further
// broken down by descriptions and price
type ReceiptImageSection struct {
	blocks  []*ocr.Block
	polygon []*ocr.Point
}

// NewReceiptImageSection creates a new image section base on a set of blocks
func NewReceiptImageSection(blocks []*ocr.Block) *ReceiptImageSection {
	if blocks == nil {
		return &ReceiptImageSection{
			blocks:  []*ocr.Block{},
			polygon: []*ocr.Point{},
		}
	}
	retval := ReceiptImageSection{
		blocks:  blocks,
		polygon: ocr.PolygonFromBlocks(blocks),
	}

	return &retval
}

// ReceiptImage is a collection of receipt image sections
type ReceiptImage struct {
	HeaderSection   *ReceiptImageSection
	LineItemSection *ReceiptImageSection
	PriceSection    *ReceiptImageSection
	SummarySection  *ReceiptImageSection
}

func (ri *ReceiptImage) String() string {
	buffer := strings.Builder{}
	buffer.WriteString("---------------------\n")
	buffer.WriteString("--- Receipt Image ---\n")
	buffer.WriteString("---------------------\n")
	buffer.WriteString("--- Header text ---\n")
	for _, block := range ri.HeaderSection.blocks {
		buffer.WriteString(fmt.Sprintf("%s ", block.Text))
	}
	buffer.WriteString("\n--- Line text ---\n")
	for _, block := range ri.LineItemSection.blocks {
		buffer.WriteString(fmt.Sprintf("%s ", block.Text))
	}
	buffer.WriteString("\n--- Price text ---\n")
	for _, block := range ri.PriceSection.blocks {
		buffer.WriteString(fmt.Sprintf("%s ", block.Text))
	}
	buffer.WriteString("\n--- Summary text ---\n")
	for _, block := range ri.SummarySection.blocks {
		buffer.WriteString(fmt.Sprintf("%s ", block.Text))
	}

	buffer.WriteString("\n")
	return buffer.String()

}

// BlockFilter is a way to encapsulate filtering logic on a block. This removes
// the need for repeatative logic (e.g. line and confidence > 90% and...)
type BlockFilter func(block *ocr.Block) bool

// always filter by line and confidence
func defaultBlockFilter(config *ImageReceiptParseConfig) BlockFilter {
	return func(block *ocr.Block) bool {
		return block.Confidence >= config.ocrConfidence
	}
}

func aboveLinearRegressionLine(lr *ocr.LinearRegression, config *ImageReceiptParseConfig, includeIntersection bool) BlockFilter {
	return func(block *ocr.Block) bool {
		// if either bottom left or bottom right are above the line, then we're good
		bottomLeft := block.BottomLeft
		bottomRight := block.BottomRight

		meetsAreaCriteria := false
		if includeIntersection {
			area, err := ocr.PercentagePolygonCoveredByLine(lr, ocr.BlockToPointArray(block), true)
			if err == nil && area >= config.minArea {
				meetsAreaCriteria = true
			}
		}

		d := defaultBlockFilter(config)(block)
		bl := bottomLeft.Y < ocr.CalculateExpectedY(bottomLeft.X, lr)
		br := bottomRight.Y < ocr.CalculateExpectedY(bottomRight.X, lr)
		a := (includeIntersection && meetsAreaCriteria)
		return d && (bl || br || a)
	}
}

func belowLinearRegressionLine(lr *ocr.LinearRegression, config *ImageReceiptParseConfig, includeIntersection bool) BlockFilter {
	return func(block *ocr.Block) bool {
		// if either top left or top right are below the line, then we're good
		topLeft := block.TopLeft
		topRight := block.TopRight

		meetsAreaCriteria := false
		if includeIntersection {
			area, err := ocr.PercentagePolygonCoveredByLine(lr, ocr.BlockToPointArray(block), false)
			if err == nil && area >= config.minArea {
				meetsAreaCriteria = true
			}
		}

		d := defaultBlockFilter(config)(block)
		tl := topLeft.Y > ocr.CalculateExpectedY(topLeft.X, lr)
		tr := topRight.Y > ocr.CalculateExpectedY(topRight.X, lr)
		a := includeIntersection && meetsAreaCriteria

		return d && (tl || tr || a)
	}
}

// NewReceiptImage creates a receipt image instance based on the textract response
func NewReceiptImage(resp *ocr.Image, config *ImageReceiptParseConfig) *ReceiptImage {

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

	return &retval
}

func findPriceViaLinearRegression(block *ocr.Block, candidateBlocks []*ocr.Block, config *ImageReceiptParseConfig) (*ocr.Block, error) {
	// calculate regression line for
	topLine, bottomLine, err := calculateSlopesForBlock(block)
	if err != nil {
		return nil, err
	}

	// remove the block from the candidate blocks to avoid picking itself
	var itr int
	for itr = 0; itr < len(candidateBlocks); itr++ {
		if candidateBlocks[itr].ID == block.ID {
			break
		}
	}

	blocks := make([]*ocr.Block, len(candidateBlocks))
	copy(blocks, candidateBlocks)
	if itr < len(candidateBlocks) {
		copy(blocks[itr:], blocks[itr+1:]) // Shift a[i+1:] left one index.
		blocks[len(blocks)-1] = nil        // Erase last element (write zero value).
		blocks = blocks[:len(blocks)-1]    // Truncate slice.
	}

	possibleItems := ocr.FindBlocksByLinearSlope(blocks, topLine, bottomLine, config.minArea)

	if len(possibleItems) == 1 {
		possibleItem := possibleItems[0]
		if defaultBlockFilter(config)(possibleItem) && priceRegex.MatchString(possibleItem.Text) {
			return possibleItem, nil
		}
	}

	// go through blocks
	for _, possibleItem := range possibleItems {
		// if it's a match, and it's a price, then return it
		if priceRegex.MatchString(possibleItem.Text) {
			return possibleItem, nil
		}
	}

	return nil, fmt.Errorf("failed to find price using block id: %s", block.ID)
}

func populateSummary(rd *ReceiptDetail, ri *ReceiptImage, config *ImageReceiptParseConfig) error {
	// typically the format is <subtotal | tax | total> (whitespace) <total value>
	// so, we will look for the the words, then verify the current line or next line is the
	// actual value

	for _, block := range ri.SummarySection.blocks {
		if subtotalRegex.MatchString(block.Text) {
			subTotalBlock, err := findPriceViaLinearRegression(block, ri.SummarySection.blocks, config)
			if err != nil {
				println(fmt.Errorf("failed to get subtotal block: %s", err.Error()))
				continue
			} else {
				val, _ := ParsePrice(subTotalBlock.Text)
				rd.SubtotalCost = val
			}

		} else if taxRegex.MatchString(block.Text) {
			taxBlock, err := findPriceViaLinearRegression(block, ri.SummarySection.blocks, config)
			if err != nil {
				return err
			}
			val, _ := ParsePrice(taxBlock.Text)
			rd.SalesTax = val
		} else if totalRegex.MatchString(block.Text) {
			totalBlock, err := findPriceViaLinearRegression(block, ri.SummarySection.blocks, config)
			if err != nil {
				return err
			}

			val, _ := ParsePrice(totalBlock.Text)
			rd.TotalCost = val
		}
	}

	// TODO: check if total or taxes are empty

	return nil
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

func filterBlocks(blocks []*ocr.Block, filter BlockFilter) []*ocr.Block {
	retval := []*ocr.Block{}
	for _, block := range blocks {
		if filter(block) {
			retval = append(retval, block)
		}
	}

	return retval
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

func calculateSlopesForPoints(pts []*ocr.Point) (*ocr.LinearRegression, *ocr.LinearRegression, error) {

	topLeft := pts[0]
	topRight := pts[1]
	bottomRight := pts[2]
	bottomLeft := pts[3]

	topLr := ocr.NewLinearRegression(topLeft.X, topLeft.Y, topRight.X, topRight.Y)
	bottomLr := ocr.NewLinearRegression(bottomLeft.X, bottomLeft.Y, bottomRight.X, bottomRight.Y)
	return topLr, bottomLr, nil

}

func calculateSlopesForBlock(block *ocr.Block) (*ocr.LinearRegression, *ocr.LinearRegression, error) {

	topLeft := block.TopLeft
	topRight := block.TopRight
	bottomRight := block.BottomRight
	bottomLeft := block.BottomLeft

	topLr := ocr.NewLinearRegression(topLeft.X, topLeft.Y, topRight.X, topRight.Y)
	bottomLr := ocr.NewLinearRegression(bottomLeft.X, bottomLeft.Y, bottomRight.X, bottomRight.Y)
	return topLr, bottomLr, nil

}

func bestEffortLineItemBlocks(
	bottomLine *ocr.LinearRegression,
	config *ImageReceiptParseConfig,
	itemBlocks []*ocr.Block) []*ocr.Block {

	return filterBlocks(
		itemBlocks,
		aboveLinearRegressionLine(bottomLine, config, true),
	)
}

func createReceiptItems(
	itemBlocks []*ocr.Block,
	finalPriceBlocks []*ocr.Block,
	config *ImageReceiptParseConfig) ([]*ReceiptItem, error) {

	items := []*ReceiptItem{}
	var prevLine *ocr.LinearRegression

	availableItemBlocks := []*ocr.Block{}
	availableItemBlocks = append(availableItemBlocks, itemBlocks...)

	for _, priceBlock := range finalPriceBlocks {
		b := make([]string, 0)
		_, bottomLine, err := calculateSlopesForBlock(priceBlock)
		if err != nil {
			return nil, err
		}

		possibleItems := bestEffortLineItemBlocks(
			bottomLine, config, availableItemBlocks,
		)

		if len(possibleItems) > 0 {

			for _, possibleItem := range possibleItems {
				// add it to the string
				b = append(b, possibleItem.Text)

				// remove it from the available item blocks
				var itr int
				for itr = 0; itr < len(availableItemBlocks); itr++ {
					if availableItemBlocks[itr].ID == possibleItem.ID {
						break
					}
				}

				copy(availableItemBlocks[itr:], availableItemBlocks[itr+1:])           // Shift a[i+1:] left one index.
				availableItemBlocks[len(availableItemBlocks)-1] = nil                  // Erase last element (write zero value).
				availableItemBlocks = availableItemBlocks[:len(availableItemBlocks)-1] // Truncate slice.
			}

			parsedPrice, err := ParsePrice(priceBlock.Text)
			if discountRegex.MatchString(priceBlock.Text) {
				parsedPrice = -parsedPrice
			}
			if err != nil {
				return nil, fmt.Errorf("%s is not a parsable price", priceBlock.Text)
			}
			item := &ReceiptItem{
				TotalCost: parsedPrice,
				Name:      strings.Join(b, " "),
			}
			items = append(items, item)
		} else {

			regressionReport := strings.Builder{}
			if prevLine != nil {
				regressionReport.WriteString(
					fmt.Sprintf("Prev line: slope: %v, intercept: %v\n", prevLine.Slope, prevLine.Intercept))
			}
			regressionReport.WriteString(
				fmt.Sprintf("Bottom line: slope: %v, intercept: %v", bottomLine.Slope, bottomLine.Intercept))

			return nil, fmt.Errorf("unable to find item desc for price block %s\nRegression Report: %s", priceBlock.ID, regressionReport.String())
		}
	}

	return items, nil
}

// ImageReceiptParseConfig - options to configure how the textract response is processed
type ImageReceiptParseConfig struct {
	// how far off the regression line will we consider a block
	regressionTolerance float64
	// the confidence of AWS on the text
	ocrConfidence float64
	// whether to include items on the header line to the header (vs line item)
	blocksOnHeaderLineAreHeader bool
	// whether to include items on the summary line to the summary (vs the line item section)
	blocksOnSummaryLineAreSummary bool
	// min area covered to include
	minArea float64
}

// NewReceiptDetailFromReceiptImage creates a receipt detail from a receipt image
func NewReceiptDetailFromReceiptImage(ri *ReceiptImage, config *ImageReceiptParseConfig) (*ReceiptDetail, error) {

	items, err := createReceiptItems(ri.LineItemSection.blocks, ri.PriceSection.blocks, config)
	if err != nil {
		return nil, err
	}
	retval := ReceiptDetail{}

	// get order date
	// TODO: scanning through all the blocks again, perhaps more logically
	//			 having data structured
	var orderDate *time.Time

	// check the header and summary section for a timestamp
	candidateBlocks := []*ocr.Block{}
	candidateBlocks = append(candidateBlocks, ri.HeaderSection.blocks...)
	candidateBlocks = append(candidateBlocks, ri.SummarySection.blocks...)
	for _, block := range candidateBlocks {
		if dateRegex.MatchString(block.Text) {

			// FIXME: we assume EST, should be passed in as a configuration option
			loc, _ := time.LoadLocation("America/New_York")

			res := dateRegex.FindStringSubmatch(block.Text)
			orderDateStr := res[1]
			for _, dateFormat := range dateFormats {
				orderDateTmp, err := time.ParseInLocation(dateFormat, orderDateStr, loc)
				if err != nil {
					continue
				}
				orderDate = &orderDateTmp
				break
			}

			if orderDate == nil {
				// warn we failed to parse the order date
				println(fmt.Sprintf("Failed to parse %s as a date time: %v, %v", orderDateStr, block.Confidence, config.ocrConfidence))
			}
		}
	}

	if orderDate != nil {
		retval.OrderTimestamp = *orderDate
	} else {
		return nil, fmt.Errorf("failed to find an order timestamp")
	}

	err = populateSummary(&retval, ri, config)
	if err != nil {
		return nil, err
	}

	retval.Items = items

	return &retval, nil
}

// ParseImageReceipt - Try multiple combinations of tolerances to get the right final value
func ParseImageReceipt(resp *ocr.Image, expectedTotal float32, confidence float64) (*ReceiptDetail, error) {

	// receipts will have blocks on the edge of the header/summary and the line items/price. This set is to allow
	// us to try different combinations, depending on the receipt
	includeBlocksOnIntersectionToHeaderSummary := [][]bool{
		{true, true},
		{false, false},
		{true, false},
		{false, true},
	}

	for regressionTolerance := 0.0; regressionTolerance < 0.05; regressionTolerance += 0.01 {

		for _, includeIntersectionOptions := range includeBlocksOnIntersectionToHeaderSummary {
			println(fmt.Sprintf("to Header: %t, to summary: %t", includeIntersectionOptions[0], includeIntersectionOptions[1]))
			config := ImageReceiptParseConfig{
				regressionTolerance:           regressionTolerance,
				ocrConfidence:                 confidence,
				blocksOnHeaderLineAreHeader:   includeIntersectionOptions[0],
				blocksOnSummaryLineAreSummary: includeIntersectionOptions[1],
				minArea:                       0.5,
			}
			ri := NewReceiptImage(resp, &config)
			if ri == nil {
				println(fmt.Sprintf("Failed to create receipt image"))
				continue
			}

			println(ri.String())

			// now convert from receipt image to receipt detail
			retval, err := NewReceiptDetailFromReceiptImage(ri, &config)
			if err != nil {
				// TODO: if it's the missing data error, we should immediately return
				//       to avoid running it multiple time
				println(fmt.Sprintf("Failed to convert to receipt detail: %s", err.Error()))
				continue
			}

			println(retval.String())

			// sum the items and check expected total
			// note we do some weird cents checking because
			// float64 and adding gives weird results
			actualTotalCents := 0
			for _, i := range retval.Items {
				actualTotalCents += int(math.Round(float64(i.TotalCost) * 100.0))
			}

			actualTotalCents += int(retval.SalesTax * 100.0)
			expectedTotalCents := int(expectedTotal * 100.0)

			if actualTotalCents == expectedTotalCents {
				println("Success!!!!")
				return retval, nil
			}
			println(fmt.Sprintf("expected %v cents, got %v cents", expectedTotalCents, actualTotalCents))

		}
	}

	return nil, fmt.Errorf("failed to find a tolerance for this receipt")

}
