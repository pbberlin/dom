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
	val, _ := strconv.Atoi(valstr)
	paramStr := r.FormValue("param")
	param, _ := strconv.Atoi(paramStr)
	pgForKey := HtmlPage{Param: param, Val: val}
	key := datastore.NewKey(ctx, htmlPageKind, pgForKey.Key(), 0, nil)
	q = q.Filter("__key__ >", key)
	q = q.Order("__key__")
	q = q.Offset(skip).Limit(5)
	var pages []HtmlPage
	keys, err := q.GetAll(ctx, &pages)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("err %v\n", err)))
	}

	logx.Debugf(r, "found %v pages", len(pages))
	for i, p := range pages {
		ts := time.Unix(p.UnixTs, 0).Format("02 Jan 2006 15:04")
		msg := fmt.Sprintf("<a href='/delete-page?key=%v'>Del</a> <a href='/show-page?key=%v'>%2d: %08v %08v %v - %v<a>\n",
			keys[i].Encode(), keys[i].Encode(), i+skip, p.Param, p.Val, p.Url, ts)
		// logx.Debugf(r, msg)
		w.Write([]byte(msg))
	}

	if skip > 0 || len(pages) >= batchSize {
		w.Write([]byte("\n"))
	}

	if skip > 0 {
		str := ` <a href='/query-pages?skip=%v&param=%v&val=%v'> << <a> `
		w.Write([]byte(fmt.Sprintf(str, skip-batchSize, param, val)))
	}

	if len(pages) >= batchSize {
		str := ` <a href='/query-pages?skip=%v&param=%v&val=%v'> >> <a> `
		w.Write([]byte(fmt.Sprintf(str, skip+batchSize, param, val)))
	}
}

func showPage(w http.ResponseWriter, r *http.Request) {

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

func deletePage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	if !logx.IsAppengine() {
		w.Write([]byte("Data Store delete is only available on app engine"))
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

	err = datastore.Delete(ctx, k)
	if err != nil {
		w.Write([]byte("<div style='white-space: pre-wrap;font-family: monospace;'>"))
		w.Write([]byte(fmt.Sprintf("err2 %v\n", err)))
		return
	}

	msg := fmt.Sprintf("deleted key %v", k)
	logx.Debugf(r, msg)
	w.Write([]byte(msg))

}
