package clean

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zew/logx"
	"github.com/zew/util"
	"golang.org/x/net/html"
)

func WriteBytesToFilename(filename string, ptrB *bytes.Buffer) {
	Bytes2File(filename, ptrB.Bytes())
}

// Dom2File writes DOM to file
func Dom2File(fn string, node *html.Node) {
	var b bytes.Buffer
	err := html.Render(&b, node)
	util.CheckErr(err)
	Bytes2File(fn, b.Bytes())
}

// Bytes2File writes bytes; creates path if neccessary
// and logs any errors even to appengine log
func Bytes2File(fn string, b []byte) {
	var err error
	err = ioutil.WriteFile(fn, b, 0)
	if err != nil {
		err = os.MkdirAll(filepath.Dir(fn), os.ModePerm)
		if err != nil {
			logx.Printf("Directory creation failed: %v", err)
			return
		}
		err = ioutil.WriteFile(fn, b, 0)
		util.CheckErr(err)
	}
}

// BytesFromFile reads bytes and logs any
// errors even to appengine log.
func BytesFromFile(fn string) []byte {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		logx.Printf("%v", err)
	}
	return b
}
