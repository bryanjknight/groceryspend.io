package parser

import (
	"fmt"
	"regexp"
	"strings"

	"groceryspend.io/server/services/receipts"
)

type Store struct {
	Name        string
	Address     string
	PhoneNumber string
}

type LineItem struct {
	Name     string
	Total    float32
	UnitCost string
	Weight   float32
}

// collection of regexes
var phoneNumberRegex = regexp.MustCompile(`(\(\d{3}\)-)?\d{3}-\d{4}`)
var addressRegex = regexp.MustCompile(`(?i)\d+ [A-Z 0-9]+ [A-Z]{2,5}`)
var townCityZipRegex = regexp.MustCompile(`(?i)\[A-Z ]+, [A-Z]{2} \d{5}(-\d{4})?`)

// TODO: use levenshtein distance to judge whether stop word hit because of misread characters
var departmentNamesRegex = regexp.MustCompile(`(?i)(dairy|produce|froduce|meat|grocery)`)

var taxRegex = regexp.MustCompile(`(?i)(total )?tax`)
var subtotalRegex = regexp.MustCompile(`(?i)subtotal`)
var totalRegex = regexp.MustCompile(`(?i)total|balance`)

var dateRegex = regexp.MustCompile(`(\d{2}/\d{2}/\d{2,4}|\d{1,2}[A-Z]{3}\d{4})`)
var timeRegex = regexp.MustCompile(`\d{1,2}:\d{2}:\d{2}`)

var priceRegex = regexp.MustCompile(`\$?\d{0,5}\.\d{2}`)
var weightRegex = regexp.MustCompile(`(?i)(\d{1,3}\.\d{2}) (lb|1b|oz)`)

/*
 *  We know the following:
 *   * Conceptually a line item has a name, a total value, maybe a unit cost, maybe a weight
 *  We also know that if it's not by weight, it'll usually be one line
 *  We know weight might come before or after the name
 * The states are Start -> Header -> Line Item(s) -> Sub Total -> Tax -> Total -> Payment -> Footer
 */

type parseContext struct {
	// states:
	// 0 - header
	// 1 - line item
	// 2 - subtotal
	// 3 - tax
	// 4 - total
	// 5 - payment
	// 6 - footer
	state       uint
	header      strings.Builder
	currentItem strings.Builder
	items       []string
	subtotal    strings.Builder
	tax         strings.Builder
	total       strings.Builder
	payment     strings.Builder
	footer      strings.Builder
}

func newParseContext() *parseContext {
	return &parseContext{
		state:       0,
		header:      strings.Builder{},
		currentItem: strings.Builder{},
		items:       []string{},
		subtotal:    strings.Builder{},
		tax:         strings.Builder{},
		total:       strings.Builder{},
		payment:     strings.Builder{},
		footer:      strings.Builder{},
	}
}

func handleHeaderState(context *parseContext, line string) {

	// if we see a department name, then we're done with header
	if departmentNamesRegex.MatchString(line) {
		context.state = 1
		return
	}

	context.header.WriteString(fmt.Sprintf(" %s", line))
}

func handleLineItem(context *parseContext, line string) {

	// if the line is subtotal, finish this line item, then run the
	// subtotal handler
	if subtotalRegex.MatchString(line) {
		context.items = append(context.items, context.currentItem.String())
		context.currentItem.Reset()
		context.state = 2
		handleSubTotal(context, line)
		return
	}

	// if the line is tax, finialize and go to the next line
	if taxRegex.MatchString(line) {
		context.items = append(context.items, context.currentItem.String())
		context.currentItem.Reset()
		context.state = 3
		handleTax(context, line)
		return
	}

	// if the line is total, finialize and go to the next line
	if totalRegex.MatchString(line) {
		context.items = append(context.items, context.currentItem.String())
		context.currentItem.Reset()
		context.state = 4
		handleTax(context, line)
		return
	}

	// for now, just write it all as one big line item
	context.currentItem.WriteString(fmt.Sprintf(" %s", line))
}

func handleSubTotal(context *parseContext, line string) {

	// if the line has a tax, then run the tax handler
	if taxRegex.MatchString(line) {
		context.state = 3
		handleTax(context, line)
		return
	}

	// if this is the price, then it's the last line of this
	// section
	if priceRegex.MatchString(line) {
		context.state = 3
	}
	context.subtotal.WriteString(fmt.Sprintf(" %s", line))

}

func handleTax(context *parseContext, line string) {
	// if we get the price, then we're at the last part of
	// the tax section
	if priceRegex.MatchString(line) {
		context.state = 4
	}
	context.tax.WriteString(fmt.Sprintf(" %s", line))
}

func handleTotal(context *parseContext, line string) {
	// if we see the price, this is probably the end of the total
	// section, so move to the next state
	if priceRegex.MatchString(line) {
		context.state = 5
	}
	context.total.WriteString(fmt.Sprintf(" %s", line))

}
func handlePayment(context *parseContext, line string) {
	// if we see a time stamp, let's assume we're at the end of the
	// payment section
	if timeRegex.MatchString(line) {
		context.state = 6
	}

	context.payment.WriteString(fmt.Sprintf(" %s", line))
}

func handleFooter(context *parseContext, line string) {
	// simple add whatever's left to the footer
	context.footer.WriteString(fmt.Sprintf(" %s", line))
}

// RegexParser uses regexes to parse a receipt's text into a structured object
func RegexParser(text string) (*receipts.ReceiptDetail, error) {

	context := newParseContext()

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		println(fmt.Sprintf("%v: %s", i, line))
		switch context.state {
		case 0:
			handleHeaderState(context, line)
		case 1:
			handleLineItem(context, line)
		case 2:
			handleSubTotal(context, line)
		case 3:
			handleTax(context, line)
		case 4:
			handleTotal(context, line)
		case 5:
			handlePayment(context, line)
		case 6:
			handleFooter(context, line)
		}
	}

	// print context for debugging purposes
	print(fmt.Sprintf(`
	---- Header ----
	%s

	---- Line Items ----
	%s

	---- Subtotal ----
	%s

	---- Tax -----
	%s

	---- Total ----
	%s

	---- Payment ----
	%s

	---- Footer ----
	%s
	`, context.header.String(),
		strings.Join(context.items, " "),
		context.subtotal.String(),
		context.tax.String(),
		context.total.String(),
		context.payment.String(),
		context.footer.String()))

	return &receipts.ReceiptDetail{}, nil

}
