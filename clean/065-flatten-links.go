package clean

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pbberlin/dom/node"
	"github.com/zew/util"

	"golang.org/x/net/html"
)

// Func funcNodeTreeToString returns a func.
// Link tags may contain images, text, and multiple nested divs.
// Func assembles such tags into a single text string.
func funcNodeTreeToString() func(n *html.Node, depth int) (b *bytes.Buffer) {
	var fLinkTreeToString func(*html.Node, int) *bytes.Buffer
	fLinkTreeToString = func(n *html.Node, depth int) (b *bytes.Buffer) {
		b = new(bytes.Buffer)
		if n.Type == html.TextNode {
			fmt.Fprint(b, strings.TrimSpace(n.Data))
		}
		fmt.Fprint(b, " ")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			b.Write(fLinkTreeToString(c, depth+1).Bytes())
		}
		return
	}
	return fLinkTreeToString
}

func flattenLinks(n *html.Node) {

	linkTree2String := funcNodeTreeToString()

	var fRecurse func(*html.Node)
	fRecurse = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "a" {

			origHref := node.AttrX(n.Attr, "href")
			hasName := (node.AttrX(n.Attr, "name") != "")
			hasHref := (origHref != "")
			startsWithRune := strings.HasPrefix(origHref, "#")

			//
			isLocalAnchor := false
			if (!hasHref || startsWithRune) && hasName {
				isLocalAnchor = true
			}

			// This is the counter part to
			// dropRedundantLinkTitles
			title := node.AttrX(n.Attr, "title")
			completeText := linkTree2String(n, 0).String()
			if title != "" {
				if strings.Contains(title, completeText) {
					completeText = title
					// logx.Printf("222: %v contained in \n%v", completeText, title)
				} else {
					completeText = fmt.Sprintf("%v. %v", title, completeText)
				}
				n.Attr = node.RemoveAttr(n.Attr, "title")
			}

			replNd := node.NewTextNode(completeText)

			// Missing link text
			// Maybe link was styled by css background-image
			if util.IsSpacey(replNd.Data) {
				if hasHref && !startsWithRune {
					// missing link text
					replNd.Data = fmt.Sprintf("[mlt] %v", util.RemoveProtocol(origHref))
				}
			}

			// Remove all existing children.
			// Direct loop impossible, since "NextSibling" is set to nil by Remove().
			// Therefore first assembling separately, then removing.
			children := make(map[*html.Node]struct{})
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				children[c] = struct{}{}
			}
			for k, _ := range children {
				n.RemoveChild(k)
			}

			//
			// Put our single textual replacement
			// under the original node
			if isLocalAnchor {
				replNd.Data = "[LocalAnchor] " + replNd.Data
				// n.Parent.InsertBefore(replNd,n)
				node.InsertAfter(n, replNd)
				node.RemoveNode(n)
			} else {
				n.AppendChild(replNd)
				// Enforce at least one space before each link
				sep1 := node.NewTextNode("  ")
				// sep2 := node.NewTextNode("\n[a]")
				// set3 := node.NewElementNode("br")
				n.Parent.InsertBefore(sep1, n)
			}

		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fRecurse(c)
		}
	}

	fRecurse(n)

}
