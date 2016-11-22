package clean

import (
	"fmt"
	"strings"

	"github.com/pbberlin/dom/node"
	"golang.org/x/net/html"
)

func dropRedundantLinkTitles(n *html.Node) {

	linkTree2String := funcNodeTreeToString()

	var fRecurse func(*html.Node)
	fRecurse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			title := node.AttrX(n.Attr, "title")
			if title != "" {
				completeText := linkTree2String(n, 0).String()
				if strings.Contains(completeText, title) {
					n.Attr = node.RemoveAttr(n.Attr, "title")
					// logx.Printf("111: %v contained in \n%v", title, completeText)
				} else {
					// We might go other way around:
					// If title contains completeText,
					// drop completeText in favor of title.
					// We keep the title for later on.
					// See 222:
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fRecurse(c)
		}
	}
	fRecurse(n)

}

func splitMultiSourceImages(n *html.Node) {

	var fRecurse func(*html.Node)
	fRecurse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for i := 0; i < 4; i++ {
				srcI := node.AttrX(n.Attr, fmt.Sprintf("src%v", i))
				if srcI != "" {
					imgI := node.NewElementNode("img")
					imgI.Attr = []html.Attribute{html.Attribute{Key: "src", Val: srcI}}
					node.InsertAfter(n, imgI)
					n.Attr = node.RemoveAttr(n.Attr, fmt.Sprintf("src%v", i))
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fRecurse(c)
		}
	}
	fRecurse(n)

}
