package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"groceryspend.io/server/services/receipts"
)

// collection of regexes
var phoneNumberRegex = regexp.MustCompile(`(\([0-9]{3}\) |[0-9]{3}-)[0-9]{3}-[0-9]{4}`)
var addressRegex = regexp.MustCompile(`(?i)\d+ [A-Z 0-9]+ [A-Z]{2,5}`)
var townCityZipRegex = regexp.MustCompile(`(?i)[A-Z ]+, [A-Z]{2} \d{5}(-\d{4})?`)
var cashierRegex = regexp.MustCompile(`(?i)cashier:? [a-z0-9 ]+`)
var storeRegex = regexp.MustCompile(`(?i)store:? \d{1,5}`)

// TODO: use levenshtein distance to judge whether stop word hit because of misread characters
var departmentNamesRegex = regexp.MustCompile(`(?i)(dairy|produce|froduce|meat|grocery)`)

var taxRegex = regexp.MustCompile(`(?i)(total )?tax`)
var subtotalRegex = regexp.MustCompile(`(?i)subtotal`)
var totalRegex = regexp.MustCompile(`(?i)total|balance`)

var dateRegex = regexp.MustCompile(`(\d{2}/\d{2}/\d{2,4}|\d{1,2}[A-Z]{3}\d{4})`)
var timeRegex = regexp.MustCompile(`\d{1,2}:\d{2}:\d{2}`)

var priceRegexStr = `(-?\d{0,5}\.\d{2})`
var priceRegex = regexp.MustCompile(priceRegexStr)

var weightRegexStr = `(?i)(\d{1,3}\.\d{2}) (lb|1b|oz)`
var weightRegex = regexp.MustCompile(weightRegexStr)

type lineItemParser struct {
	name              string
	regex             *regexp.Regexp
	qtyGroupID        uint
	weightGroupID     uint
	weightUnitGroupID uint
	unitPriceGroupID  uint
	nameGroupID       uint
	finalPriceGroupID uint
}

// these are applied in order; thus, the most complex should go first
var lineItemParsers = []*lineItemParser{
	// <weight> @ <unit price> <name> <final price>
	// ex: 0.55 lb @ 1 1b / 2.99 GR/HOUSE RED PEPPERS 1.64 F
	{
		name:              "Weight-UnitPrice-Name-FinalPrice",
		regex:             regexp.MustCompile(`^(?i)((\d{1,3}\.\d{2}) (lb|1b|oz))( @ )?(\d (lb|1b) \/ \d{1,2}\.\d{2}) (([A-Z0-9 /#\&]){3,}) (\d{0,3}\.\d{2})( [a-z\*] ?)?`),
		weightGroupID:     2,
		weightUnitGroupID: 3,
		unitPriceGroupID:  5,
		nameGroupID:       7,
		finalPriceGroupID: 9,
	},
	// <name> <weight> @ <unit price> <final price>
	// ex: HOT HOUSE TOMATOES W 0.63 lb @ 2.99/ 1b 1.88 *
	{
		name:              "Name-Weight-UnitPrice-FinalPrice",
		regex:             regexp.MustCompile(`^(?i)(([A-Z0-9 /#\&]){3,}) (\d{1,3}\.\d{2}) (lb|1b|oz)( @ )?(\d{0,2}\.\d{2}(\/ 1b|lb|oz)) (\d{0,3}\.\d{2})( [a-z\*] ?)?`),
		nameGroupID:       1,
		weightGroupID:     3,
		weightUnitGroupID: 4,
		unitPriceGroupID:  6,
		finalPriceGroupID: 8,
	},
	// <qty> @ <unit price> <name> <final price>
	// ex: 1 @ 2/ 5.00 DRAG WHL MOZZ CHUNK 2.50 F
	{
		name:              "Qty-UnitPrice-Name-FinalPrice",
		regex:             regexp.MustCompile(`^(?i)(\d{1,3})( @ )?(\d ?\/ \d{1,2}\.\d{2}) ([A-Z0-9 \/#\&]{3,}) (\d{0,3}\.\d{2})( [a-z\*] ?)?`),
		qtyGroupID:        1,
		unitPriceGroupID:  3,
		nameGroupID:       4,
		finalPriceGroupID: 5,
	},
	// <name> <qty> @ <unit price> <final price>
	// ex: BLACKBERRIES W 2 @ 3.49 6.98 *
	{
		name:              "Name-Qty-UnitPrice-FinalPrice",
		regex:             regexp.MustCompile(`^(?i)(([A-Z0-9 /#\&]){3,}) (\d{1,3})( @ )?(\d{0,2}\.\d{2}) (\d{0,3}\.\d{2})( [a-z\*] ?)?`),
		nameGroupID:       1,
		qtyGroupID:        3,
		unitPriceGroupID:  5,
		finalPriceGroupID: 6,
	},
	// <name> <discount>
	// ex: MB RESTAURANT TORTS 2.00 F
	{
		name:              "Name-Discount",
		regex:             regexp.MustCompile(`^(?i)(([A-Z0-9 /#\&]){3,}) (-\d{0,3}\.\d{2})( [a-z\*] ?)?`),
		nameGroupID:       1,
		finalPriceGroupID: 3,
	},
	// <name> <price>
	// ex: MB RESTAURANT TORTS 2.00 F
	{
		name:              "Name-FinalPrice",
		regex:             regexp.MustCompile(`^(?i)(([A-Z0-9 /#\&]){3,}) (\d{0,3}\.\d{2})( [a-z\*] ?)?`),
		nameGroupID:       1,
		finalPriceGroupID: 3,
	},
}

// parseState is an enum of the various parsing states
type parseState int

// Valid parseState
const (
	Header parseState = iota + 1
	LineItems
	Subtotal
	Tax
	Total
	Payment
	Footer
)

// String - Creating common behavior - give the type a String function
func (d parseState) String() string {
	return [...]string{"Header", "LineItems", "Subtotal", "Tax", "Total", "Payment", "Footer"}[d-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (d parseState) EnumIndex() int {
	return int(d)
}

type parseContext struct {
	state          parseState
	header         strings.Builder
	itemBuffer     strings.Builder
	items          []*receipts.ReceiptItem
	subtotalBuffer strings.Builder
	subtotal       float32
	taxBuffer      strings.Builder
	tax            float32
	totalBuffer    strings.Builder
	total          float32
	payment        strings.Builder
	footer         strings.Builder
}

func newParseContext() *parseContext {
	return &parseContext{
		state:          Header,
		header:         strings.Builder{},
		itemBuffer:     strings.Builder{},
		items:          []*receipts.ReceiptItem{},
		subtotalBuffer: strings.Builder{},
		taxBuffer:      strings.Builder{},
		totalBuffer:    strings.Builder{},
		payment:        strings.Builder{},
		footer:         strings.Builder{},
	}
}

func handleHeaderState(context *parseContext, line string) error {

	// if we see a department name, then we're done with header
	if departmentNamesRegex.MatchString(line) {
		context.state = LineItems
	}

	context.header.WriteString(fmt.Sprintf(" %s", line))
	return nil
}

func finalizeTax(context *parseContext) error {
	// get the price and set it as the tax
	res := priceRegex.FindStringSubmatch(context.taxBuffer.String())
	if res == nil {
		println("didn't find a tax value, assuming not supplied")
		return nil
	}

	tax, err := strconv.ParseFloat(res[1], 32)
	if err != nil {
		return err
	}
	context.tax = float32(tax)
	return nil
}

func finalizeSubtotal(context *parseContext) error {
	// get the price and set it as the subtotal
	res := priceRegex.FindStringSubmatch(context.subtotalBuffer.String())
	if res == nil {
		println("didn't find a subtotal value, assuming not supplied")
		return nil
	}

	subtotal, err := strconv.ParseFloat(res[1], 32)
	if err != nil {
		return err
	}
	context.subtotal = float32(subtotal)
	return nil
}

func finalizeTotal(context *parseContext) error {
	// get the price and set it as the subtotal
	res := priceRegex.FindStringSubmatch(context.totalBuffer.String())
	if res == nil {
		return fmt.Errorf("didn't find a total value")
	}

	total, err := strconv.ParseFloat(res[1], 32)
	if err != nil {
		return err
	}
	context.total = float32(total)
	return nil
}

func finalizeLineItems(context *parseContext) error {
	// we will try each line item parser until we get a match. We'll then parse out the details (unit price, total cost, name, etc)
	// and return it as a ReceiptItem
	itemsStr := strings.TrimSpace(context.itemBuffer.String())
	context.itemBuffer.Reset()

	for len(itemsStr) > 0 {
		isProcessed := false
		for _, lineItemParser := range lineItemParsers {
			// if we have already processed the line, start from the beginning
			if isProcessed {
				break
			}
			res := lineItemParser.regex.FindStringSubmatch(itemsStr)
			if res == nil {
				continue
			}
			name := res[lineItemParser.nameGroupID]
			tmp, err := strconv.ParseFloat(res[lineItemParser.finalPriceGroupID], 32)
			totalPrice := float32(tmp)
			if err != nil {
				return err
			}
			qty := 0
			if lineItemParser.qtyGroupID != 0 {
				qty, err = strconv.Atoi(res[lineItemParser.qtyGroupID])
				if err != nil {
					println(err)
					qty = 0
				}
			}

			weight := float32(0.0)
			if lineItemParser.weightGroupID != 0 {
				tmp, err = strconv.ParseFloat(res[lineItemParser.weightGroupID], 32)
				weight = float32(tmp)
				if err != nil {
					println(err)
					weight = 0
				}
			}
			unit := ""
			if lineItemParser.weightUnitGroupID != 0 {
				unit = res[lineItemParser.weightUnitGroupID]
			}

			// unitPrice := ""
			// if lineItemParser.unitPriceGroupID != 0 {
			// 	unitPrice = res[lineItemParser.unitPriceGroupID]
			// }

			context.items = append(context.items, &receipts.ReceiptItem{
				Name:          name,
				TotalCost:     totalPrice,
				Weight:        float32(weight),
				ContainerUnit: unit,
				// TODO: parse unit cost UnitCost:      unitPrice,
				Qty: qty,
			})

			itemsStr = itemsStr[len(res[0]):]
			isProcessed = true
		}

		if !isProcessed {
			println(itemsStr)
			return fmt.Errorf("failed to match any parser to line (see stdout)")
		}
	}

	return nil
}

func handleLineItem(context *parseContext, line string) error {

	// if the line is a deparatment, skip
	if departmentNamesRegex.MatchString(line) {
		return nil
	}

	// if the line is subtotal, finish this line item, then run the
	// subtotal handler
	if subtotalRegex.MatchString(line) {
		err := finalizeLineItems(context)
		if err != nil {
			return err
		}
		context.state = Subtotal
		handleSubTotal(context, line)
		return nil
	}

	// if the line is tax, finialize and go to the next line
	if taxRegex.MatchString(line) {
		err := finalizeLineItems(context)
		if err != nil {
			return err
		}
		context.state = Tax
		handleTax(context, line)
		return nil
	}

	// if the line is total, finialize and go to the next line
	if totalRegex.MatchString(line) {
		err := finalizeLineItems(context)
		if err != nil {
			return err
		}
		context.state = Total
		handleTax(context, line)
		return nil
	}

	// for now, just write it all as one big line item
	context.itemBuffer.WriteString(fmt.Sprintf(" %s", line))
	return nil
}

func handleSubTotal(context *parseContext, line string) error {

	// if the line has a tax, then run the tax handler
	if taxRegex.MatchString(line) {
		context.state = Tax
		handleTax(context, line)
		return nil
	}

	// if this is the price, then it's the last line of this
	// section
	if priceRegex.MatchString(line) {
		context.subtotalBuffer.WriteString(fmt.Sprintf(" %s", line))
		err := finalizeSubtotal(context)
		if err != nil {
			return err
		}
		context.state = Tax
	} else {
		context.subtotalBuffer.WriteString(fmt.Sprintf(" %s", line))
	}
	return nil
}

func handleTax(context *parseContext, line string) error {
	// if we get the price, then we're at the last part of
	// the tax section
	if priceRegex.MatchString(line) {
		context.taxBuffer.WriteString(fmt.Sprintf(" %s", line))
		finalizeTax(context)
		context.state = Total
	} else {
		context.taxBuffer.WriteString(fmt.Sprintf(" %s", line))
	}
	return nil
}

func handleTotal(context *parseContext, line string) error {
	// if we see the price, this is probably the end of the total
	// section, so move to the next state
	if priceRegex.MatchString(line) {
		context.totalBuffer.WriteString(fmt.Sprintf(" %s", line))
		err := finalizeTotal(context)
		if err != nil {
			return err
		}
		context.state = Payment
	} else {
		context.totalBuffer.WriteString(fmt.Sprintf(" %s", line))
	}

	return nil

}
func handlePayment(context *parseContext, line string) error {
	// if we see a time stamp, let's assume we're at the end of the
	// payment section
	if timeRegex.MatchString(line) {
		context.state = Footer
	}

	context.payment.WriteString(fmt.Sprintf(" %s", line))
	return nil
}

func handleFooter(context *parseContext, line string) error {
	// simple add whatever's left to the footer
	context.footer.WriteString(fmt.Sprintf(" %s", line))
	return nil
}

// RegexParser uses regexes to parse a receipt's text into a structured object
func RegexParser(text string) (*receipts.ReceiptDetail, error) {

	context := newParseContext()

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		// println(fmt.Sprintf("%v: %s", i, line))
		var err error
		switch context.state {
		case parseState(Header.EnumIndex()):
			err = handleHeaderState(context, line)
		case parseState(LineItems.EnumIndex()):
			err = handleLineItem(context, line)
		case parseState(Subtotal.EnumIndex()):
			err = handleSubTotal(context, line)
		case parseState(Tax.EnumIndex()):
			err = handleTax(context, line)
		case parseState(Total.EnumIndex()):
			err = handleTotal(context, line)
		case parseState(Payment.EnumIndex()):
			err = handlePayment(context, line)
		case parseState(Footer.EnumIndex()):
			err = handleFooter(context, line)
		}

		if err != nil {
			return nil, err
		}
	}

	return &receipts.ReceiptDetail{
		Items:        context.items,
		SubtotalCost: context.subtotal,
		TotalCost:    context.total,
		SalesTax:     context.tax,
	}, nil

}
