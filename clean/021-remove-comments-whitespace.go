package clean

import (
	"github.com/pbberlin/dom/node"
	"github.com/zew/util"
	"golang.org/x/net/html"
)

// removeCommentsAndInterTagWhitespace employs horizontal traversal using a queue
func removeCommentsAndInterTagWhitespace(lp interface{}) {

	var queue = NewQueue(10)

	for lp != nil {

		lpn := lp.(node.NdX).Nd
		lvl := lp.(node.NdX).Lvl

		for c := lpn.FirstChild; c != nil; c = c.NextSibling {
			queue.EnQueue(node.NdX{c, lvl + 1})
		}

		// processing
		if lpn.Type == html.CommentNode {
			node.RemoveNode(lpn)
		}

		// extinguish textnodes that do only formatting (spaces, tabs, line breaks)
		if lpn.Type == html.TextNode && util.IsSpacey(lpn.Data) {
			node.RemoveNode(lpn)
		}

		// next node
		lp = queue.DeQueue()
	}
}
