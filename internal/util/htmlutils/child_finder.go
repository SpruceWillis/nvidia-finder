package htmlutils

import (
	"golang.org/x/net/html"
)

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

// FindChildByAttribute find a child node with a certain attribute - use nil if unset
func FindChildByAttribute(node *html.Node, targetAttr html.Attribute) *html.Node {
	condition := func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if attr.Namespace == targetAttr.Namespace && attr.Key == targetAttr.Key && attr.Val == targetAttr.Val {
				return true
			}
		}
		return false
	}
	return FindChild(node, condition)
}
