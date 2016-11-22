package clean

import (
	"github.com/pbberlin/dom/node"
	"github.com/zew/logx"
	"golang.org/x/net/html"
)

var debugBreakOut = false

const breakoutMarker = " [boi=>]"         // broken out image: Next element
const breakoutMarkerSelf = " [boi-lts=>]" // broken out image: link to self

func searchImg(n *html.Node, fnd *html.Node, lvl int) (*html.Node, int) {

	if n.Type == html.ElementNode && n.Data == "img" {
		// logx.Printf("  a has img on lvl %v\n", lvl)
		if fnd == nil {
			fnd = n
			return fnd, lvl
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		fnd, lvlfnd := searchImg(c, fnd, lvl+1)
		if fnd != nil {
			return fnd, lvlfnd
		}
	}

	return fnd, lvl
}

type DeleterFunc func(*html.Node, int, bool) bool

func closureDeleter(until bool) DeleterFunc {

	// Nodes along the path to the splitting image
	// should never be removed in *neither* tree
	var splitPath = map[*html.Node]bool{}

	var fc DeleterFunc
	fc = func(n *html.Node, lvl int, found bool) bool {

		// fmt.Printf("found %v at l%v\n", found, lvl)
		if n.Data == "img" {
			// fmt.Printf(" found at l%v\n", lvl)
			found = true
			par := n.Parent
			for {
				if par == nil {
					break
				}
				splitPath[par] = true
				par = par.Parent
			}
		}

		// children
		cc := []*html.Node{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cc = append(cc, c)
		}
		for _, c := range cc {
			found = fc(c, lvl+1, found)
		}

		//
		// remove
		if lvl > 0 {
			if n.Data == "img" {
				n.Parent.RemoveChild(n)
			} else {
				if !until && !found && !splitPath[n] {
					n.Parent.RemoveChild(n)
				}
				if until && found && !splitPath[n] {
					n.Parent.RemoveChild(n)
				}
			}
		}

		return found

	}

	return fc

}

// Searches images inside link elements.
// Cuts those images into two halves.
// Puts the img in between.
// Img is *not yet* converted to link.
func breakoutImgFromLinkTrees(n *html.Node) {

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		breakoutImgFromLinkTrees(c)
	}

	if n.Type == html.ElementNode && n.Data == "a" {

		img, lvl := searchImg(n, nil, 0)

		if img != nil {

			only1Child := n.FirstChild != nil && n.FirstChild == n.LastChild

			// Link has only one child: an image.
			// => Replace with text node.
			// => Put the image next to the link.
			if lvl == 1 && only1Child {
				// logx.Printf("only child image lvl %v a\n", lvl)

				// Special case: The link leads to the img src:
				linkToImgItself := false
				if node.AttrX(n.Attr, "href") == node.AttrX(img.Attr, "src") {
					linkToImgItself = true
				}

				n.RemoveChild(img)
				n.Parent.InsertBefore(img, n.NextSibling) // "insert after; if n.NextSibling==nil => insert at the end"

				attrDigest := node.AttributeDigest(n)

				if linkToImgItself {
					node.RemoveNode(n)
					node.InsertBefore(img, node.NewTextNode(breakoutMarkerSelf)) // broken out image next
				} else {
					n.AppendChild(node.NewTextNode(attrDigest)) // inject a linktext for the removed image
					// Put some marker behind the link - before the img
					node.InsertAfter(n, node.NewTextNode(breakoutMarker)) // broken out image next
				}

			} else {

				if debugBreakOut {
					b0 := node.PrintSubtree(n)
					logx.Printf("\n%s\n", b0)
				}
				// logx.Printf("  got it  %v\n", img.Data)
				a1 := node.CloneNodeWithSubtree(n)
				fc1 := closureDeleter(true)
				fc1(n, 0, false)
				if debugBreakOut {
					b1 := node.PrintSubtree(n)
					logx.Printf("\n%s\n", b1)
				}

				fc2 := closureDeleter(false)
				fc2(a1, 0, false)
				if debugBreakOut {
					b2 := node.PrintSubtree(a1)
					logx.Printf("\n%s\n", b2)
					logx.Printf("--------------------\n")
				}

				n.Parent.InsertBefore(img, n.NextSibling) // "insert after; if n.NextSibling==nil => insert at the end"
				n.Parent.InsertBefore(a1, img.NextSibling)

				attrDigest := "[attrdigest] " + node.AttributeDigest(n)
				attrDigest = "" // not neccessary

				node.InsertBefore(img, node.NewTextNode(attrDigest+breakoutMarker)) // broken out image next

			}

			// changing image to link later

		} else {
			// logx.Printf("no img in a\n")
		}
	}

}
