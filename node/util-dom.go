// Package node contains dom manipulations;
// inspired by github.com/PuerkitoBio/goquery/blob/master/manipulation.go
package node

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/zew/logx"
	"github.com/zew/util"
	"golang.org/x/net/html"
)

// NdX is a html.node, extended by its level.
// It's used since the horizontal traversal with
// a queue has no recursion and therefore
// keeps no depth information.
type NdX struct {
	Nd  *html.Node
	Lvl int
}

type NodeTypeStr html.NodeType

func (n NodeTypeStr) String() string {
	switch n {
	case 0:
		return "ErroNd"
	case 1:
		return "Text  "
	case 2:
		return "DocmNd"
	case 3:
		return "Elem  "
	case 4:
		return "CommNd"
	case 5:
		return "DoctNd"
	}
	return "unknown Node type"
}

// Create a new node
func NewTextNode(content string) *html.Node {
	nd0 := new(html.Node)
	nd0.Type = html.TextNode
	nd0.Data = content
	return nd0
}

func NewElementNode(kind string) *html.Node {
	nd0 := new(html.Node)
	nd0.Type = html.ElementNode
	nd0.Data = kind
	return nd0
}

func unused__ReplaceNode(self, dst *html.Node) {
	InsertAfter(self, dst)
	RemoveNode(self)
}

func RemoveNode(n *html.Node) {
	par := n.Parent
	if par != nil {
		par.RemoveChild(n)
	} else {
		logx.Printf("\nNode to remove has no Parent\n")
		logx.PrintStackTrace()
	}
}

// InsertBefore inserts before itself.
// node.InsertBefore refers to its children
func InsertBefore(insPnt, toInsert *html.Node) {
	if insPnt.Parent != nil {
		insPnt.Parent.InsertBefore(toInsert, insPnt)
	} else {
		logx.Printf("\nInsertBefore - insPnt has no Parent\n")
		logx.PrintStackTrace()
	}
}

// InsertBefore inserts at the end, when NextSibling is null.
// compare http://stackoverflow.com/questions/4793604/how-to-do-insert-after-in-javascript-without-using-a-library
func InsertAfter(insPnt, toInsert *html.Node) {
	if insPnt.Parent != nil {
		insPnt.Parent.InsertBefore(toInsert, insPnt.NextSibling)
	} else {
		logx.Printf("\nInsertAfter - insPnt has no Parent\n")
		logx.PrintStackTrace()
	}
}

//
//
// Deep copy a node.
// The new node has clones of all the original node's
// children but none of its parents or siblings
func CloneNodeWithSubtree(n *html.Node) *html.Node {
	nn := &html.Node{
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]html.Attribute, len(n.Attr)),
	}

	copy(nn.Attr, n.Attr)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nn.AppendChild(CloneNodeWithSubtree(c)) // recursion
	}
	return nn
}

//
//
// Deep copy a node.
// no children, no parent, no siblings
func CloneNode(n *html.Node) *html.Node {
	nn := &html.Node{
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]html.Attribute, len(n.Attr)),
	}
	copy(nn.Attr, n.Attr)
	return nn
}

//
func PrintSubtree(n *html.Node) *bytes.Buffer {
	b := new(bytes.Buffer)
	return printSubtree(n, b, 0)
}
func printSubtree(n *html.Node, b *bytes.Buffer, lvl int) *bytes.Buffer {

	if lvl > 40 {
		logx.Printf("%s", b.String())
		logx.Printf("possible circular relationship\n")
		os.Exit(1)
	}

	ind := strings.Repeat("  ", lvl)
	slvl := fmt.Sprintf("%sL%v", ind, lvl)
	fmt.Fprintf(b, "%-10v %v", slvl, NodeTypeStr(n.Type))
	fmt.Fprintf(b, " %v\n", strings.TrimSpace(n.Data))

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		printSubtree(c, b, lvl+1) // recursion
	}

	return b
}

func AttributeDigest(n *html.Node) string {

	elements := make([]string, 1, 2)

	// alt := AttrX(n.Attr, "alt") // alt and title are unified before
	title := AttrX(n.Attr, "title")

	if title != "" {
		elements[0] = title
	} else {
		src := AttrX(n.Attr, "src")
		if src == "" {
			src = AttrX(n.Attr, "href")
		}
		elements[0] = util.UrlBeautify(src)
	}
	// superset := util.TrimRedundant(elements)
	return elements[0]

}

// Beware: There can be *multiple* identical keys
// Breaking after the first finding
// would return the first one.
func AttrX(attrs []html.Attribute, key string) string {
	if key == "src" {
		src := attrX(attrs, key)
		if src == "" {
			if srcAlt := attrX(attrs, "srcset"); srcAlt != "" {
				src = srcAlt
			}
		}
		if src == "" {
			if srcAlt := attrX(attrs, "data-src"); srcAlt != "" {
				src = srcAlt
			}
		}
		return src
	}

	return attrX(attrs, key)
}
func attrX(attrs []html.Attribute, key string) string {
	ret := ""
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Key == key {
			if len(attrs[i].Val) > len(ret) { // longest value wins
				ret = attrs[i].Val
			}
		}
	}
	return ret
}

func AttrSet(attrs []html.Attribute, key, val string) []html.Attribute {
	for i, a := range attrs {
		if a.Key == key {
			attrs[i].Val = val
			return attrs
		}
	}
	// attr does not exist => append it
	attrs = append(attrs, html.Attribute{Key: key, Val: val})
	return attrs
}

func Unused_addIdAttr(attributes []html.Attribute, id string) []html.Attribute {
	hasId := false
	for _, a := range attributes {
		if a.Key == "id" {
			hasId = true
			break
		}
	}
	if !hasId {
		attributes = append(attributes, html.Attribute{"", "id", id})
	}
	return attributes
}

// Primitive performant attr removal.
func RemoveAttr(attrs []html.Attribute, removeKey string) []html.Attribute {
	if removeKey == "src" {
		attrs = removeAttr(attrs, "data-src")
		attrs = removeAttr(attrs, "srcset")
	}
	return removeAttr(attrs, removeKey)
}
func removeAttr(attrs []html.Attribute, removeKey string) []html.Attribute {
	exists := 0
	for i := 0; i < len(attrs); i++ {
		attrs[i].Key = strings.TrimSpace(attrs[i].Key)
		attrs[i].Key = strings.ToLower(attrs[i].Key)
		if attrs[i].Key == removeKey {
			exists++
		}
	}

	if exists < 1 {
		return attrs
	}

	repl := make([]html.Attribute, 0, len(attrs)-exists)

	for i := 0; i < len(attrs); i++ {
		if attrs[i].Key != removeKey {
			repl = append(repl, attrs[i])
		}
	}
	return repl
}

// bool true => remove prefix
// bool false => remove exact match
func RemoveAttrPrefix(attrs []html.Attribute, removeKeys map[string]bool) []html.Attribute {

	ret := []html.Attribute{}

NextAttr:
	for _, a := range attrs {
		a.Key = strings.TrimSpace(a.Key)
		a.Key = strings.ToLower(a.Key)
		a.Val = strings.TrimSpace(a.Val)
		a.Val = util.NormalizeInnerWhitespace(a.Val) // having encountered title or alt values with newlines

		// Keep these for those cases,
		// where ONLY src rewrite is applied
		// so that javascripts can use this values
		// to set the "onload" src
		if a.Key == "data-src" || a.Key == "srcset" || a.Key == "data-href" {
			ret = append(ret, a)
			continue NextAttr
		}

		for unwanted, alsoPrefix := range removeKeys {
			if alsoPrefix {
				if strings.HasPrefix(a.Key, unwanted) {
					continue NextAttr // throw away
				}
			} else {
				if a.Key == unwanted {
					continue NextAttr // throw away
				}
			}
		}
		ret = append(ret, a)
	}
	return ret
}
