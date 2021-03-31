package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type InstacartTipTaxesFees struct {
	SalesTax    float32
	Tip         float32
	ServiceFee  float32
	DeliveryFee float32
	Discounts   float32
}

func parseLineItem(li *html.Node) (ParsedItem, error) {
	// this is where assumptions are made, and thus the most likely
	// part to fail. We ASSUME that each line item is wrapped in one div, which then has two
	// separate divs: one for the description the second for the qty/weight and the final price

	// thoughts -- we know the data is stored in p tags, so why not just get those
	// using our parser and get the underlying data from the TextElement

	topLevelDiv := li.FirstChild

	descriptionDiv := topLevelDiv.FirstChild
	descriptionPTags := GetElementsByTagName(descriptionDiv, "p")

	if len(descriptionPTags) < 2 {
		return ParsedItem{}, errors.New("failed to retrieve description tags")
	}

	name := descriptionPTags[0].FirstChild.Data
	unitPriceAndContainerSize := descriptionPTags[1].FirstChild.Data
	tmpSlice := strings.Split(unitPriceAndContainerSize, "•")
	unitPrice, err := ParseStringToUSDAmount(tmpSlice[0])
	if err != nil {
		return ParsedItem{}, errors.New("failed to parse Unit Price")
	}
	// containerSize := ""
	// if len(tmpSlice) > 1 {
	// 	containerSize = tmpSlice[1]
	// }

	qtyDiv := topLevelDiv.LastChild
	qtyPTags := GetElementsByTagName(qtyDiv, "p")

	if len(qtyPTags) < 2 {
		return ParsedItem{}, errors.New("failed to retrieve qty tags")
	}

	qtyString := qtyPTags[0].FirstChild.Data
	qty, err := strconv.Atoi(qtyString)
	if err != nil {
		qty = 0
	}
	totalPrice, err := ParseStringToUSDAmount(qtyPTags[1].FirstChild.Data)
	if err != nil {
		return ParsedItem{}, errors.New("failed to parse total price")
	}

	retval := ParsedItem{}
	retval.Name = name
	retval.Qty = qty
	retval.UnitCost = unitPrice
	retval.TotalCost = totalPrice
	// retval.ContainerSize = containerSize
	return retval, nil
}

func parseItemsFound(ul *html.Node) ([]ParsedItem, error) {

	retval := []ParsedItem{}

	// iterate through each item
	for li := ul.FirstChild; li != nil; li = li.NextSibling {
		pi, err := parseLineItem(li)
		if err != nil {
			println(err)
			continue
		}
		retval = append(retval, pi)
	}
	return retval, nil
}

func parseReplacementsAndRefunded(children []*html.Node) ([]ParsedItem, error) {
	retval := []ParsedItem{}

	// if there's more than 2, then we have either
	// replacements
	// refunded
	// replacements and refunded
	// we need to look at the 3rd and (if necessary) the 5th elements to see

	// if they're "replacements" and "refunded"
	if len(children) == 4 {
		// determine if the second section is refunded or replacements
		h3Tags := GetElementsByTagName(children[2], "h3")
		if len(h3Tags) != 1 {
			println("failed to determine section type, skipping")
			return retval, nil
		}

		text := h3Tags[0].FirstChild.Data

		if strings.Contains(text, "replacement") {
			// item 4 in the array is the replacements
			replacements, err := parseItemsFound(children[3])
			if err != nil {
				return retval, err
			}
			retval = append(retval, replacements...)
		}

	} else if len(children) == 6 {

		// item 4 in the array is the replacements
		replacements, err := parseItemsFound(children[3])
		if err != nil {
			return retval, err
		}
		retval = append(retval, replacements...)

	}
	return retval, nil
}

func parseTaxTipFees(sectionDiv *html.Node) (InstacartTipTaxesFees, error) {

	retval := InstacartTipTaxesFees{}

	divs := GetElementsByTagName(sectionDiv, "div")

	for _, div := range divs {

		pTags := GetElementsByTagName(div, "p")

		if len(pTags) != 2 {
			println("Didn't get two p tags for an element, skipping")
			continue
		}

		name := pTags[0].FirstChild.Data
		cost, err := ParseStringToUSDAmount(pTags[1].FirstChild.Data)
		if err != nil {
			println("Unable to parse %v as a money value, skipping", pTags[1].FirstChild.Data)
			continue

		}

		switch name {
		case "Sales Tax":
			retval.SalesTax = cost
		case "Tip":
			retval.Tip = cost
		case "Delivery Fee":
			retval.DeliveryFee = cost
		case "Service Fee":
			retval.ServiceFee = cost
		case "Deals Discount":
			retval.Discounts = cost
		case "Item Subtotal":
			// no op
		case "Total":
			// no op
		default:
			println(fmt.Sprintf("Unexpected subtotal value: %v", name))
		}
	}

	return retval, nil
}

func ParseInstcartHtmlReceipt(doc *html.Node) (ParsedReceipt, error) {
	// find the "main" tag
	mainNodes := GetElementsByTagName(doc, "main")
	if len(mainNodes) != 1 {
		return ParsedReceipt{}, fmt.Errorf("expected 1 main node, got %v", len(mainNodes))
	}

	mainNode := mainNodes[0]

	// making this an array makes it easier to work with
	children := []*html.Node{}
	for c := mainNode.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	// find the found items
	itemsFound, err := parseItemsFound(children[1])
	if err != nil {
		// this is not recoverable, so it implies a parsing issue
		return ParsedReceipt{}, err
	}

	replacements, err := parseReplacementsAndRefunded(children)
	if err != nil {
		// this is not recoverable, so it implies a parsing issue
		return ParsedReceipt{}, err
	}

	itemsFound = append(itemsFound, replacements...)

	// get the taxes, tips, and feeds
	// it's main's parent's next sibling's child, so main's cousin ¯\_(ツ)_/¯
	taxTipFeesDiv := mainNode.Parent.NextSibling.FirstChild
	taxTipFees, err := parseTaxTipFees(taxTipFeesDiv)
	if err != nil {
		return ParsedReceipt{}, err
	}

	retval := ParsedReceipt{}
	retval.ParsedItems = itemsFound
	retval.DeliveryFee = taxTipFees.DeliveryFee
	retval.SalesTax = taxTipFees.SalesTax
	retval.ServiceFee = taxTipFees.ServiceFee
	retval.Discounts = taxTipFees.Discounts
	retval.Tip = taxTipFees.Tip

	return retval, nil
}
