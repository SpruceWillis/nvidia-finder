package htmlutils

import "golang.org/x/net/html"

// GetAttribute return a node's attribute, if it exists
func GetAttribute(node *html.Node, attribute string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == attribute {
			return attr.Val, true
		}
	}
	return "", false
}
