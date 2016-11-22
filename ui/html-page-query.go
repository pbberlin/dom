package ui

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zew/logx"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const batchSize = 5

//
func queryPages(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.Write([]byte("<div style='white-space: pre-wrap;font-family: monospace;'>"))
	if !logx.IsAppengine() {
		w.Write([]byte("Data Store query is only available on app engine"))
		return
	}
	ctx := appengine.NewContext(r)

	skipstr := r.FormValue("skip")
	skip, _ := strconv.Atoi(skipstr)

	q := datastore.NewQuery(htmlPageKind)
	valstr := r.FormValue("val")
	val, err := strconv.Atoi(valstr)
	if err != nil {
	} else {
		q = q.Filter("Val >=", val)
		logx.Debugf(r, "val is %v", val)
	}
	q = q.Order("Val").Order("Url").Order("UnixTs")
	q = q.Offset(skip).Limit(5)
	var pages []HtmlPage
	keys, err := q.GetAll(ctx, &pages)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v\n", err)))
	}

	logx.Debugf(r, "found %v pages", len(pages))
	for i, p := range pages {
		ts := time.Unix(p.UnixTs, 0).Format("02 Jan 2006 15:04")
		msg := fmt.Sprintf("<a href='/show-detail?key=%v'>%2d: %08v %v - %v<a>\n", keys[i].Encode(), i+skip, p.Val, p.Url, ts)
		// logx.Debugf(r, msg)
		w.Write([]byte(msg))
	}

	if skip > 0 || len(pages) >= batchSize {
		w.Write([]byte("\n"))
	}

	if skip > 0 {
		str := ` <a href='/query-pages?skip=%v&val=%v'> << <a> `
		w.Write([]byte(fmt.Sprintf(str, skip-batchSize, val)))
	}

	if len(pages) >= batchSize {
		str := ` <a href='/query-pages?skip=%v&val=%v'> >> <a> `
		w.Write([]byte(fmt.Sprintf(str, skip+batchSize, val)))
	}
}

func showDetail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	if !logx.IsAppengine() {
		w.Write([]byte("Data Store detail query is only available on app engine"))
		return
	}
	ctx := appengine.NewContext(r)

	keystr := r.FormValue("key")
	k, err := datastore.DecodeKey(keystr)
	if err != nil {
		w.Write([]byte("<div style='white-space: pre-wrap;font-family: monospace;'>"))
		w.Write([]byte(fmt.Sprintf("err1 %v \n\t%v\n", err, k)))
		return
	}

	logx.Debugf(r, "looking up key %+v", k)

	pg := &HtmlPage{}
	err = datastore.Get(ctx, k, pg)
	if err != nil {
		w.Write([]byte("<div style='white-space: pre-wrap;font-family: monospace;'>"))
		w.Write([]byte(fmt.Sprintf("err2 %v\n", err)))
		return
	}

	logx.Debugf(r, "found result for url %v val %v - %v bytes", pg.Url, pg.Val, len(pg.Body))
	// logx.Debugf(r, "%v", util.Ellipsoider(string(pg.Body), 50))

	// w.Write([]byte(util.IndentedDump(pg)))
	w.Write(pg.Body)

}
