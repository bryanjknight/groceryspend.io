package receipts

import (
	"fmt"
	"math"
	"strings"
	"time"

	"groceryspend.io/server/services/ocr"
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

	if rd.TotalCost == 0.0 {
		return fmt.Errorf("total cost is zero")
	}

	return nil
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

	// We go through increasingly larger regression tolerances to find the best fit. In experiments, using a static
	// tolerance value would work for some receipts and fail on other receipts.
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
			ri, err := NewReceiptImage(resp, expectedTotal, &config)
			if err != nil {
				println(err.Error())
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
