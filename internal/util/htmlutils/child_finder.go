package htmlutils

import "golang.org/x/net/html"

// FindChild find a child node that meets a condition. Returns nil if not found
func FindChild(node *html.Node, childConditions func(*html.Node) bool) *html.Node {
	child := node.FirstChild
	for child != nil {
		if childConditions(child) {
			return child
		}
		child = child.NextSibling
	}
	return nil
}
