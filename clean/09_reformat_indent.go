package clean

import (
	"log"
	"strings"

	"github.com/pbberlin/dom/node"
	"golang.org/x/net/html"
)

func ReIndent(n *html.Node, lvl int) {

	if lvl > cScaffoldLvls && n.Parent == nil {
		bb := node.PrintSubtree(n)
		_ = bb
		// logx.Printf("%s", bb.Bytes())
		hint := ""
		if ml3[n] > 0 {
			hint = "   from ml3"
		}
		log.Print("reIndent: no parent ", hint)
		return
	}

	// Before children processing
	switch n.Type {
	case html.ElementNode:
		if lvl > cScaffoldLvls && n.Parent.Type == html.ElementNode {
			ind := strings.Repeat("\t", lvl-2)
			node.InsertBefore(n, &html.Node{Type: html.TextNode, Data: "\n" + ind})
		}
	case html.CommentNode:
		node.InsertBefore(n, &html.Node{Type: html.TextNode, Data: "\n"})
	case html.TextNode:
		n.Data = strings.TrimSpace(n.Data) + " "
		if !strings.HasPrefix(n.Data, ",") && !strings.HasPrefix(n.Data, ".") {
			n.Data = " " + n.Data
		}
		// link texts without trailing space
		if n.Parent != nil && n.Parent.Data == "a" {
			n.Data = strings.TrimSpace(n.Data)
		}
	}

	// Children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ReIndent(c, lvl+1)
	}

	// After children processing
	switch n.Type {
	case html.ElementNode:
		// I dont know why,
		// but this needs to happend AFTER the children
		if lvl > cScaffoldLvls && n.Parent.Type == html.ElementNode {
			ind := strings.Repeat("\t", lvl-2)
			ind = "\n" + ind
			// link texts without new line
			if n.Data == "a" {
				ind = ""
			}
			if n.LastChild != nil {
				node.InsertAfter(n.LastChild, &html.Node{Type: html.TextNode, Data: ind})
			}
		}
	}

}
