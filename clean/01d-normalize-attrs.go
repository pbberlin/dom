package clean

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var sources = map[string]string{
	"data-src": "src",
	"srcset":   "src",
}
var harefs = map[string]string{
	"data-href": "href",
	"onclick":   "href",
}

func normalizeAttrs(attrs []html.Attribute) []html.Attribute {
	rew := make([]html.Attribute, 0, len(attrs)) // rewritten
	srcs := []string{"default"}                  // first entry for src
	hrefs := []string{"default"}                 // first entry for href
	alttitle := []string{}                       // title and or alt
	for i := 0; i < len(attrs); i++ {
		attrs[i].Val = strings.TrimSpace(attrs[i].Val)
		attrs[i].Key = strings.TrimSpace(attrs[i].Key)
		attrs[i].Key = strings.ToLower(attrs[i].Key)
		if attrs[i].Key == "src" {
			srcs[0] = attrs[i].Val
			continue
		}
		if _, ok := sources[attrs[i].Key]; ok { // alternative srcs
			srcs = append(srcs, attrs[i].Val)
			continue
		}
		if attrs[i].Key == "style" { // extract bg images - otherwise continue
			if bgUrl := extractBackgroundImgUrl(attrs[i].Val); bgUrl != "" {
				srcs = append(srcs, bgUrl)
			}
			continue
		}
		if attrs[i].Key == "href" {
			hrefs[0] = attrs[i].Val
			continue
		}
		if _, ok := harefs[attrs[i].Key]; ok { // alternative hrefs
			hrefs = append(hrefs, attrs[i].Val)
			continue
		}
		if attrs[i].Key == "alt" || attrs[i].Key == "title" {
			alttitle = append(alttitle, attrs[i].Val)
			continue
		}
		rew = append(rew, attrs[i])
	}
	rew = appendSrcHref(rew, srcs, "src")
	rew = appendSrcHref(rew, hrefs, "href")
	if len(alttitle) == 1 {
		rew = append(rew, html.Attribute{Key: "title", Val: alttitle[0]})
	}
	if len(alttitle) > 1 {
		if strings.Contains(alttitle[0], alttitle[1]) {
			rew = append(rew, html.Attribute{Key: "title", Val: alttitle[0]})
		} else if strings.Contains(alttitle[1], alttitle[0]) {
			rew = append(rew, html.Attribute{Key: "title", Val: alttitle[1]})
		} else {
			rew = append(rew, html.Attribute{Key: "title", Val: strings.Join(alttitle, ". ")})
		}
	}
	return rew
}

// A slice of prioritized href or src values
// is appended to attrs.
func appendSrcHref(attrs []html.Attribute, vals []string, key string) []html.Attribute {
	if len(vals) == 1 && vals[0] != "default" {
		attrs = append(attrs, html.Attribute{Key: key, Val: vals[0]})
	}
	if len(vals) > 1 {
		idxStart := 0
		if vals[0] == "default" {
			idxStart = 1
		}
		cntr := 0
		for i := idxStart; i < len(vals); i++ {
			if strings.Contains(vals[i], "/") {
				if cntr == 0 {
					attrs = append(attrs, html.Attribute{Key: key, Val: vals[i]})
				} else {
					attrs = append(attrs, html.Attribute{Key: fmt.Sprintf("%v%v", key, i), Val: vals[i]})
				}
				cntr++
			}
		}
	}
	return attrs
}

// xxx-url(/img1.jpg)--xxx
//   =>
// /img1.jpg
const bgImgUrlPref = "url("
const lenUrlBrack = len(bgImgUrlPref)

func extractBackgroundImgUrl(style string) string {
	if pos1 := strings.Index(style, bgImgUrlPref); pos1 > -1 {
		pos1 = pos1 + lenUrlBrack
		offs := strings.Index(style[pos1:], ")")
		if offs > -1 {
			style = style[pos1 : pos1+offs]
			style = strings.Trim(style, "\"' ")
			return style
		}
	}
	return ""
}

func throwAwayAttrs(attrs []html.Attribute) []html.Attribute {
	rew := make([]html.Attribute, 0, len(attrs))
	for i := 0; i < len(attrs); i++ {
		attr := attrs[i]
		switch attr.Key {
		case "cfrm":
			// keep
		case "href", "action", "src", "encoding":
			// keep
		case "src0", "src1", "src2", "src3":
			// keep
		case "href0", "href1", "href2", "href3":
			// keep
		case "type", "name", "value", "accesskey":
			// keep
		case "title":
			// keep
		default:
			if !strings.Contains(attr.Val, "/") {
				continue
			}
			// throw away
		}
		rew = append(rew, attr)
	}
	return rew

}

func normalizeAttrsTree(n *html.Node) {

	var rec func(n *html.Node, lvl int)
	rec = func(n *html.Node, lvl int) {
		for c := n.FirstChild; c != nil; c = c.NextSibling { // Children
			rec(c, lvl+1)
		}
		n.Attr = normalizeAttrs(n.Attr)
	}
	rec(n, 0)

}

func throwAwayAttrsTree(n *html.Node) {

	var rec func(n *html.Node, lvl int)
	rec = func(n *html.Node, lvl int) {
		for c := n.FirstChild; c != nil; c = c.NextSibling { // Children
			rec(c, lvl+1)
		}
		n.Attr = throwAwayAttrs(n.Attr)
	}
	rec(n, 0)

}
