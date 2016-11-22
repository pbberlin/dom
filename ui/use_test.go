package ui

import (
	"net/http"
	"testing"

	"github.com/pbberlin/dom/clean"
	"github.com/zew/logx"
)

// Instead of go build && run,
// use   go test
func Test_Build(t *testing.T) {

	opt1 := func(c *clean.Config) { c.HtmlTitle = "Proxify http requests" }
	cf.Apply(opt1, opt1)

	ExplicitInit(nil)

	logx.Printf("\nstarting server on 4072")
	logx.Fatal(http.ListenAndServe("localhost:4072", nil))
}
