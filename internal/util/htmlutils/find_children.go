package htmlutils

import (
	"strings"

	"golang.org/x/net/html"
)

// FindChildren find all direct descendants of a node that meet a certain condition
func FindChildren(node *html.Node, f func(*html.Node) bool) []*html.Node {
	nodes := make([]*html.Node, 0)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if f(child) {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// FindChildrenByClass find all direct descendants with a class
func FindChildrenByClass(node *html.Node, exactMatch bool, targetClass string) []*html.Node {
	condition := func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if attr.Key == class {
				if exactMatch {
					if attr.Val == targetClass {
						return true
					}
				} else {
					strings := strings.Split(attr.Val, " ")
					for _, s := range strings {
						if s == targetClass {
							return true
						}
					}
				}
			}
		}
		return false
	}
	return FindChildren(node, condition)
}
