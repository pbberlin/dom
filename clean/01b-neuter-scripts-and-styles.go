package clean

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pbberlin/dom/node"
	"github.com/zew/logx"
	"github.com/zew/util"
	"golang.org/x/net/html"
)

const emptySrc = "//:0"

func neuterScriptsAndStyles(n *html.Node) {

	tagCounter := 0 // for closure neuterNode

	// Reducing nodes to a textnode summary of their contents
	// We could replace these with a link (<a>),
	// but some these occur in html head, where <a> links are
	// impossible
	neuterNode := func(n *html.Node) {

		if n.Type != html.ElementNode {
			return
		}

		if n.Data != "script" &&
			n.Data != "noscript" &&
			n.Data != "style" &&
			n.Data != "link" &&
			n.Data != "iframe" &&
			n.Data != "object" &&
			n.Data != "canvas" &&
			n.Data != "svg" &&
			n.Data != "picture" &&
			n.Data != "map" {
			return
		}

		var replNd *html.Node
		b := new(bytes.Buffer) // contents of single text childnode
		if n.Data == "script" {
			src := node.AttrX(n.Attr, "src")
			n.Attr = []html.Attribute{} // clear all attributes
			fmt.Fprintf(b, " var script%02v = '[script]'; // src %v", tagCounter, src)
			n.Attr = node.AttrSet(n.Attr, "src", emptySrc)
			tagCounter++
		}
		if n.Data == "noscript" {
			fmt.Fprint(b, "[noscript_contents_removed]")
		}
		if n.Data == "style" {
			href := node.AttrX(n.Attr, "href")
			n.Attr = []html.Attribute{}
			fmt.Fprintf(b, " .dummyclass {margin:2px;} /* href was %v */", href)
			// tagCounter++
		}
		if n.Data == "link" {
			href := node.AttrX(n.Attr, "href")
			n.Attr = []html.Attribute{html.Attribute{Key: "removed-href", Val: href}}
			// tagCounter++
			return // done with link
		}
		if n.Data == "iframe" {
			n.Data = "a" // convert to link
			href := node.AttrX(n.Attr, "src")
			n.Attr = []html.Attribute{html.Attribute{Key: "href", Val: href}}
			n.Attr = append(n.Attr, html.Attribute{Key: "cfrm", Val: "iframe"})
			fmt.Fprintf(b, "[iframe_removed] %v", util.UrlBeautify(href))
			// tagCounter++
		}
		if n.Data == "object" {
			n.Data = "a" // convert to link
			href := node.AttrX(n.Attr, "data")
			stype := node.AttrX(n.Attr, "type")
			n.Attr = []html.Attribute{html.Attribute{Key: "href", Val: href}}
			n.Attr = append(n.Attr, html.Attribute{Key: "cfrm", Val: "object"})
			fmt.Fprintf(b, "[object_removed] %v type -%v-", util.UrlBeautify(href), stype)
			// tagCounter++
		}
		if n.Data == "canvas" {
			n.Data = "div" // convert to div
			id := node.AttrX(n.Attr, "id")
			n.Attr = append(n.Attr, html.Attribute{Key: "cfrm", Val: "canvas"})
			fmt.Fprintf(b, "[canvas_removed] id -%v-", id)
			// tagCounter++
		}
		if n.Data == "svg" {
			n.Data = "div" // convert to div
			n.Attr = append(n.Attr, html.Attribute{Key: "cfrm", Val: "svg"})
			svgTree2String := funcNodeTreeToString()
			completeText := svgTree2String(n, 0).String()
			completeText = strings.TrimSpace(completeText)
			if completeText != "" {
				fmt.Fprintf(b, "[svg_removed] cnt -%v-", completeText)
			} else {
				fmt.Fprintf(b, "[svg_removed]")
			}
			// tagCounter++
		}
		if n.Data == "picture" {
			n.Data = "div" // convert to div
			n.Attr = append(n.Attr, html.Attribute{Key: "cfrm", Val: "picture"})
			replNd = node.NewElementNode("div")
			srcCounter := 0
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				logx.Printf("pic child %v", c.Data)
				if c.Data == "source" {
					src := node.AttrX(c.Attr, "srcset")
					slSrc := strings.Fields(src)
					if len(slSrc) > 0 {
						// fmt.Fprintf(b, " src%v -%v- ", srcCounter+1, slSrc[0])
						imgSrc := node.NewElementNode("img")
						imgSrc.Attr = []html.Attribute{html.Attribute{Key: "src", Val: slSrc[0]}}
						replNd.AppendChild(imgSrc)
					}
					srcCounter++
				}
				if srcCounter > 1 {
					break
				}
			}
		}
		if n.Data == "map" {
			n.Data = "div" // convert to div
			n.Attr = append(n.Attr, html.Attribute{Key: "cfrm", Val: "map"})
			fmt.Fprintf(b, "[map_removed]")
			// tagCounter++
		}

		if replNd == nil {
			replNd = node.NewTextNode(b.String())
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
		// Put our single textual replacement
		// under the original node
		n.AppendChild(replNd)
	}

	var rec func(n *html.Node, lvl int)
	rec = func(n *html.Node, lvl int) {
		// Children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rec(c, lvl+1)
		}
		neuterNode(n)
	}
	rec(n, 0)

}
