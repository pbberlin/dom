package clean

import (
	"github.com/zew/util"
	"golang.org/x/net/html"
)

// !DOCTYPE html head
// !DOCTYPE html body
//        0    1    2
const cScaffoldLvls = 2

var (
	ml3 = map[*html.Node]int{}

	nodeDistinct = map[string]int{}
	attrDistinct = map[string]int{}

	// thrown away completely
	unwanteds = map[string]bool{
		"script":   true,
		"noscript": true,
		// "style":    true,
		// "link":     true,
		// "iframe":   true,
		// "object":   true,
		// "canvas":   true,

		"wbr": true,

		"meta":    true,
		"comment": true,
	}

	skip = map[string]bool{
		"br": true,
	}

	// converted to div
	// note that
	// 		<table><tr><td>
	// 		</td></tr></table>
	// is not changed
	exotics = map[string]string{
		"header":  "div",
		"footer":  "div",
		"nav":     "div",
		"section": "div",
		"article": "div",
		"aside":   "div",

		"fieldset": "div", // check this

		"dl": "ul",
		"dt": "li",
		"dd": "p",

		"figure":     "div",
		"figcaption": "p",

		"i": "em",
		"b": "strong",
	}
)

// maxTreeDepth returns the depth of given DOM node
func maxTreeDepth(n *html.Node, lvl int) (maxLvl int) {
	maxLvl = lvl
	for c := n.FirstChild; c != nil; c = c.NextSibling { // Children
		ret := maxTreeDepth(c, lvl+1)
		if ret > maxLvl {
			maxLvl = ret
		}
	}
	return
}

// convertUnwanted neutralizes a node.
// Note: We can not directly Remove() nor Replace()
// Since that breaks the recursion one step above!
// At a later stage we employ horizontal traversal
// to actually remove unwanted nodes.
//
// Meanwhile we have devised removeUnwanted() which
// makes convertUnwanted-removeComment obsolete.
//
func convertUnwanted(n *html.Node) {
	if unwanteds[n.Data] {
		n.Type = html.CommentNode
		n.Data = n.Data + " replaced"
	}
}

// We want to remove some children.
// A direct loop is impossible,
// since "NextSibling" is set to nil during Remove().
// Therefore:
//   First assemble children separately.
//   Then remove them.
func removeUnwanted(n *html.Node) {
	cc := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cc = append(cc, c)
	}
	for _, c := range cc {
		if unwanteds[c.Data] {
			n.RemoveChild(c)
		}
	}
}

// convertExotic standardizes <section> or <header> nodes
// towards <div> nodes.
func convertExotic(n *html.Node) {
	if repl, ok := exotics[n.Data]; ok {
		n.Attr = append(n.Attr, html.Attribute{"", "cfrm", n.Data})
		n.Data = repl
	}
}

// purgeScriptsAndStyles performs brute reduction and simplification
//
func purgeScriptsAndStyles(n *html.Node, lvl int) {

	for c := n.FirstChild; c != nil; c = c.NextSibling { // Children
		purgeScriptsAndStyles(c, lvl+1)
	}

	if true {
		removeUnwanted(n) // direct removal now working
	} else {
		convertUnwanted(n)
	}
	convertExotic(n)

	// one time text normalization
	if n.Type == html.TextNode {
		n.Data = util.NormalizeInnerWhitespace(n.Data)
	}

}
