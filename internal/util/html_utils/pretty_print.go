package htmlutils

import (
	"fmt"

	"golang.org/x/net/html"
)

// PrettyPrintHTMLNode print useful information about the node and its children to a given depth
func PrettyPrintHTMLNode(node *html.Node, depth int) {
	if depth <= 0 {
		return
	}
	fmt.Println("node:", node)
	switch node.Type {
	case html.ErrorNode:
		fmt.Println("error node")
	case html.TextNode:
		fmt.Println("text node")
	case html.DocumentNode:
		fmt.Println("document node")
	case html.ElementNode:
		fmt.Println("document node")
	case html.CommentNode:
		fmt.Println("comment node")
	case html.DoctypeNode:
		fmt.Println("doctype node")
	case html.RawNode:
		fmt.Println("raw node")
	}
	numSiblings, numChildren := 0, 0
	sibling := node.NextSibling
	for sibling != nil {
		numSiblings++
		sibling = sibling.NextSibling
	}
	fmt.Printf("node has %v siblings", numSiblings)
	child := node.FirstChild
	for child != nil {
		numChildren++
		child = child.NextSibling
	}
	fmt.Printf("node has %v children", numChildren)
	if depth > 1 {
		child = node.FirstChild
		for child != nil {
			PrettyPrintHTMLNode(child, depth-1)
			child = child.NextSibling
		}
	}
}
