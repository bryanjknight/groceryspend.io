package receipts

import (
	"regexp"

	"golang.org/x/net/html"
)

func checkTagName(n *html.Node, tagName string) bool {
	if n.Type == html.ElementNode {
		if n.Data == tagName {
			return true
		}
	}
	return false
}

func traverse(n *html.Node, test func(*html.Node) bool) []*html.Node {
	retval := []*html.Node{}

	if test(n) {
		retval = append(retval, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, test)
		if len(result) > 0 {
			retval = append(retval, result...)
		}
	}

	return retval
}

// GetElementsByTagName recurse through DOM to find nodes that match this tag name. Does not include self-closing tags
func GetElementsByTagName(node *html.Node, tagName string) []*html.Node {

	testFunc := func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			if n.Data == tagName {
				return true
			}
		}
		return false
	}
	return traverse(node, testFunc)
}

// GetAttribute retrieve a node's attribute
func GetAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkID(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := GetAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}
	return false
}

func traverseByID(n *html.Node, id string) *html.Node {
	if checkID(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverseByID(c, id)
		if result != nil {
			return result
		}
	}

	return nil
}

// GetElementByID find element by it's HTML ID
func GetElementByID(n *html.Node, id string) *html.Node {
	return traverseByID(n, id)
}

// GetElementByTextContent find the node with the specific tag name and has a value that matches the regular expression
func GetElementByTextContent(node *html.Node, tagName string, re *regexp.Regexp) []*html.Node {
	testFunc := func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == tagName {

			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode && re.MatchString(n.FirstChild.Data) {
				return true
			}
			return false

		}
		return false
	}

	return traverse(node, testFunc)
}
