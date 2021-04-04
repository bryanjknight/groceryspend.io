package receipts

import (
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

func traverseByTagName(n *html.Node, name string) []*html.Node {
	retval := []*html.Node{}

	if checkTagName(n, name) {
		retval = append(retval, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverseByTagName(c, name)
		if len(result) > 0 {
			retval = append(retval, result...)
		}
	}

	return retval
}

// GetNodesByTagName recurse through DOM to find nodes that match this tag name.
// 										does not include self-closing tags
func GetElementsByTagName(node *html.Node, tagName string) []*html.Node {
	return traverseByTagName(node, tagName)
}

func GetAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := GetAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}
	return false
}

func traverseById(n *html.Node, id string) *html.Node {
	if checkId(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverseById(c, id)
		if result != nil {
			return result
		}
	}

	return nil
}

func GetElementById(n *html.Node, id string) *html.Node {
	return traverseById(n, id)
}
