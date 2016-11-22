package clean

import (
	"fmt"
	"strings"

	"github.com/pbberlin/dom/node"
	"github.com/zew/util"

	"golang.org/x/net/html"
)

// Is there already a text node
// in the vicinity of img
// containing txt
func closureTextNodeExists(img *html.Node, txt string) (found bool) {

	if len(txt) < 5 {
		return false
	}
	txt = util.NormalizeInnerWhitespace(txt)
	txt = strings.TrimSpace(txt)

	// We dont search entire document, but three levels above image subtree
	grandParent := img
	for i := 0; i < 4; i++ {
		if grandParent.Parent != nil {
			grandParent = grandParent.Parent
		} else {
			// logx.Printf("LevelsUp %v for %q", i, txt)
			break
		}
	}

	var recurseTextNodes func(n *html.Node)
	recurseTextNodes = func(n *html.Node) {

		if found {
			return
		}

		cc := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cc = append(cc, c)
		}
		for _, c := range cc {
			recurseTextNodes(c)
		}

		if n.Type == html.TextNode {
			n.Data = util.NormalizeInnerWhitespace(n.Data)
			if len(n.Data) >= len(txt) {
				// if strings.Contains(txt, "FDP") {
				// 	logx.Printf("%25v     %v", util.Ellipsoider(txt, 10), util.Ellipsoider(n.Data, 10))
				// }
				fnd := strings.Contains(n.Data, txt)
				if fnd {
					found = true
					return
				}
			}
		}
	}
	recurseTextNodes(grandParent)

	return
}

func nodeToLink(img *html.Node) {

	if img.Data == "img" {

		// Convert to link - and convert src to href
		img.Data = "a"

		src := node.AttrX(img.Attr, "src")
		img.Attr = node.AttrSet(img.Attr, "href", src)
		img.Attr = node.AttrSet(img.Attr, "src", src)
		// We dont delete src for performance reasons

		img.Attr = node.RemoveAttr(img.Attr, "src")

		// if !strings.Contains(src, "gravatar") {
		// 	logx.Printf("changed img %v", src)
		// }

		img.Attr = node.AttrSet(img.Attr, "href", src)
		img.Attr = node.AttrSet(img.Attr, "cfrom", "img")

		attrDigest := node.AttributeDigest(img)
		double := closureTextNodeExists(img, attrDigest)
		if double {
			attrDigest = "[ctdr]" // content title double removed
		}

		s := fmt.Sprintf("[img] %v", attrDigest)
		linkText := node.NewTextNode(s)

		// Since img is now a link "a", it can contain child nodes:
		img.AppendChild(linkText)
		node.InsertBefore(img, node.NewTextNode(" "))    // put a space *before* the link
		node.InsertAfter(img, node.NewTextNode(" "))     // put a space *before* the link
		node.InsertAfter(img, node.NewElementNode("br")) // put a break *after*  the link

		return

		// It seemd we could not put our replacement
		// under an image converted to link.

		//
		replacementLink := node.NewElementNode("a")
		replacementLink.Attr = node.AttrSet(img.Attr, "href", src)
		replacementLink.Attr = node.AttrSet(img.Attr, "cfrom", "img")
		replacementLink.AppendChild(linkText)

		// img.Parent.InsertBefore(replNd,img)
		node.InsertAfter(img, replacementLink)
		node.RemoveNode(img)

		// Now we need to do this:
		(*img) = *replacementLink
		// Otherwise, the change is only reflected in the DOM.
		// But in anchorFlattener, we delete the DOM,
		// and read the replNd instead.

	}

}

func img2Link(n *html.Node) {

	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		img2Link(c)
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		nodeToLink(n)
	}
}
