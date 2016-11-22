package clean

import (
	"github.com/pbberlin/dom/node"
	"golang.org/x/net/html"
)

// Especially form action, anchor href and image src.
func closuredProxifier(p *Proxifier) func(*html.Node) {

	// Recurse points form actions back to the proxy.
	// Recurse rewrites href and src attributes point to the proxy.
	var fRecurse func(*html.Node)
	fRecurse = func(n *html.Node) {

		switch {
		// Form stuff
		case n.Type == html.ElementNode && n.Data == "form":

			// rewrite action attr to local form redirector
			origAction := node.AttrX(n.Attr, "action")
			n.Attr = p.attrsAbsoluteAndProxified(n.Attr)
			n.Attr = append(n.Attr, html.Attribute{Key: "method", Val: "post"})
			n.Attr = append(n.Attr, html.Attribute{Key: "style",
				Val: "margin: 24px 1px; border: 1px solid #aaa;"})

			// add original url for form redirector
			hidFld := node.NewElementNode("input")
			hidFld.Attr = []html.Attribute{
				html.Attribute{Key: "name", Val: p.FormRedirectKey},
				html.Attribute{Key: "value", Val: p.absolutize(p.RemoteHost(), origAction, false)},
			}
			n.AppendChild(hidFld)

			// add a visible submit button
			submt := new(html.Node)
			submt.Type = html.ElementNode
			submt.Data = "input"
			submt.Attr = []html.Attribute{
				html.Attribute{Key: "type", Val: "submit"},
				html.Attribute{Key: "value", Val: "subm"},
				html.Attribute{Key: "accesskey", Val: "s"},
			}
			n.AppendChild(submt)

			n.Attr = node.AttrSet(n.Attr, "method", "post")
			n.Attr = node.AttrSet(n.Attr, "was", "rewritten")

		case n.Type == html.ElementNode && (n.Data == "a" || n.Data == "img"):

			n.Attr = p.attrsAbsoluteAndProxified(n.Attr)

		default:
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fRecurse(c)
		}
	}

	return fRecurse

}

func (p *Proxifier) Rewrite(n *html.Node) {
	fRecurser := closuredProxifier(p)
	fRecurser(n)
}
