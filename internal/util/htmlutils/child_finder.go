package htmlutils

import (
	"strings"

	"golang.org/x/net/html"
)

const class = "class"

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

// FindChildByAttribute find a child node with a certain attribute
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

// FindChildWithAttributeKey find first child node with the desired attribute key
func FindChildWithAttributeKey(node *html.Node, key string) *html.Node {
	condition := func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if attr.Key == key {
				return true
			}
		}
		return false
	}
	return FindChild(node, condition)
}

// FindChildWithClass find first child node with desired class. Specify exactMatch = true to exactly match node classes, otherwise just include class
func FindChildWithClass(node *html.Node, exactMatch bool, targetClass string) *html.Node {
	condition := func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if attr.Key == class {
				if exactMatch {
					if attr.Val == targetClass {
						return true
					}
				} else {
					classes := strings.Split(attr.Val, " ")
					for _, c := range classes {
						if c == targetClass {
							return true
						}
					}
				}
			}
		}
		return false
	}
	return FindChild(node, condition)
}

//  FindChildWithClasses find first child with all classes present in a single class declaration
func FindChildWithClasses(node *html.Node, classes []string) *html.Node {
	condition := func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if attr.Key == class {
				targetClassesFound := make(map[string]bool)
				for _, c := range classes {
					targetClassesFound[c] = false
				}
				nodeClasses := strings.Split(attr.Val, " ")
				for _, c := range nodeClasses {
					targetClassesFound[c] = true
				}
				for _, found := range targetClassesFound {
					if !found {
						return false
					}
				}
				return true // all classes found
			}
		}
		return false
	}
	return FindChild(node, condition)
}
