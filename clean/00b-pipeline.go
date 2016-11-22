package clean

import (
	"bytes"
	"path/filepath"

	"github.com/pbberlin/dom/node"
	"github.com/zew/logx"
	"golang.org/x/net/html"
)

func FileNamer(logdir string, fileNumber int) func() string {
	cntr := -2
	return func() string {
		cntr++
		if cntr == -1 {
			return spf("outp_%03v", fileNumber) // prefix/filekey
		} else {
			fn := spf("outp_%03v_%v", fileNumber, cntr) // filename with stage
			fn = filepath.Join(logdir, fn)
			return fn
		}
	}
}

func fileDump(doc *html.Node, fNamer func() string) {
	if fNamer != nil {
		removeCommentsAndInterTagWhitespace(node.NdX{doc, 0})
		ReIndent(doc, 0)
		Dom2File(fNamer()+".html", doc)
		removeCommentsAndInterTagWhitespace(node.NdX{doc, 0})
	}
}

//
func (cf *Config) ProcessFull(b []byte) (*html.Node, error) {

	b = globalByteFixes(b)
	s2 := primitiveReformat(string(b))
	b = []byte(s2)

	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		logx.Printf("%v", err)
		return nil, err
	}

	if cf.FNamer != nil {
		Dom2File(cf.FNamer()+".html", doc)
	}

	//
	//
	neuterScriptsAndStyles(doc)
	fileDump(doc, cf.FNamer)

	purgeScriptsAndStyles(doc, 0)
	normalizeAttrsTree(doc)
	throwAwayAttrsTree(doc)

	removeCommentsAndInterTagWhitespace(node.NdX{doc, 0})
	fileDump(doc, cf.FNamer)

	//
	//
	condenseTopDown(doc, 0, 0)
	removeEmptyNodes(doc, 0)
	fileDump(doc, cf.FNamer)

	//
	//
	removeCommentsAndInterTagWhitespace(node.NdX{doc, 0}) // prevent spacey textnodes around singl child images
	breakoutImgFromLinkTrees(doc)
	dropRedundantLinkTitles(doc)
	splitMultiSourceImages(doc)

	img2Link(doc)
	flattenLinks(doc)

	fileDump(doc, cf.FNamer)

	//
	//
	condenseBottomUpV3(doc, 0, 7, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 6, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 5, map[string]bool{"div": true})
	condenseBottomUpV3(doc, 0, 4, map[string]bool{"div": true})
	condenseTopDown(doc, 0, 0)

	removeEmptyNodes(doc, 0)
	removeEmptyNodes(doc, 0)

	fileDump(doc, cf.FNamer)

	//
	//
	if !cf.SkipProxify {
		cf.Rewrite(doc)
		fileDump(doc, cf.FNamer)
	}

	if cf.Beautify {
		removeCommentsAndInterTagWhitespace(node.NdX{doc, 0})
		ReIndent(doc, 0)
	}

	//
	//
	if cf.AddOutline {
		addOutlineAttr(doc, 0, []int{0})
	}
	if cf.AddID {
		addIdAttr(doc, 0, 1)
	}
	if cf.AddOutline || cf.AddID {
		fileDump(doc, cf.FNamer)
	}

	//
	computeXPathStack(doc, 0)
	if cf.FNamer != nil {
		Bytes2File(cf.FNamer()+".txt", xPathDump)
	}

	return doc, nil

}
