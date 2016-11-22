package clean

import (
	"bytes"
	"strings"

	"github.com/pbberlin/dom/node"
	"github.com/zew/util"
	"golang.org/x/net/html"
)

func DomFormat(doc *html.Node) {
	removeEmptyNodes(doc, 0)
	removeCommentsAndInterTagWhitespace(node.NdX{doc, 0})
	ReIndent(doc, 0)
}

func (cf *Config) ProcessLean(b []byte) []byte {

	b = globalByteFixes(b)
	s2 := primitiveReformat(string(b))

	doc, err := html.Parse(strings.NewReader(s2))
	util.CheckErr(err)

	neuterScriptsAndStyles(doc)

	// changed default: Do purge
	if !cf.ProcessOptions.SkipPurgeScriptsAndStyles {
		purgeScriptsAndStyles(doc, 0)
	}
	normalizeAttrsTree(doc)
	throwAwayAttrsTree(doc)

	removeCommentsAndInterTagWhitespace(node.NdX{doc, 0})
	condenseTopDown(doc, 0, 0)
	removeEmptyNodes(doc, 0)

	breakoutImgFromLinkTrees(doc)
	dropRedundantLinkTitles(doc)
	splitMultiSourceImages(doc)

	if !cf.SkipImage2Link {
		img2Link(doc)
	}
	flattenLinks(doc)

	condenseBottomUpV3(doc, 0, 7, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 6, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 5, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 4, map[string]bool{"div": true})
	condenseTopDown(doc, 0, 0)

	removeEmptyNodes(doc, 0)
	removeEmptyNodes(doc, 0)

	if !cf.ProcessOptions.SkipProxify {
		cf.Proxifier.Rewrite(doc)
	}

	DomFormat(doc)

	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	util.CheckErr(err)

	return buf.Bytes()

}
