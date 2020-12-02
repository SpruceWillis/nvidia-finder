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
	var nodeType string
	switch node.Type {
	case html.ErrorNode:
		nodeType = "error node"
	case html.TextNode:
		nodeType = "text node"
	case html.DocumentNode:
		nodeType = "document node"
	case html.ElementNode:
		nodeType = "element node"
	case html.CommentNode:
		nodeType = "comment node"
	case html.DoctypeNode:
		nodeType = "doctype node"
	case html.RawNode:
		nodeType = "raw node"
	}
	fmt.Println("dataAtom:", node.DataAtom)
	fmt.Println("node type:", nodeType)
	fmt.Println("data:", node.Data)
	fmt.Println("namespace:", node.Namespace)
	fmt.Println("attributes:", node.Attr)
	numSiblings, numChildren := 0, 0
	parent := node.Parent
	if parent != nil {
		sibling := parent.FirstChild
		for sibling != nil {
			numSiblings++
			sibling = sibling.NextSibling
		}
		numSiblings-- // do not count the node itself as a sibling
	}
	fmt.Println(fmt.Sprintf("node has %v siblings", numSiblings))
	child := node.FirstChild
	for child != nil {
		numChildren++
		child = child.NextSibling
	}
	fmt.Println(fmt.Sprintf("node has %v children", numChildren))
	if depth > 1 {
		child = node.FirstChild
		for child != nil {
			PrettyPrintHTMLNode(child, depth-1)
			child = child.NextSibling
		}
	}
}
