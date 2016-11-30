package ui

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/monoculum/formam"
	"github.com/pbberlin/dom/clean"
	"github.com/pbberlin/fetch"
	"github.com/zew/logx"
	"github.com/zew/util"
)

const valCustom = -5 // for selective deletion

type FormData struct {
	Url string `formam:"urlx" json:"urlx"` // config param url_param_key must match this :(
	// Specs    []struct {
	// 	SpecId int
	// }
	clean.ProcessOptions

	Submit string `formam:"submit2"`
}

func writeHeader(w http.ResponseWriter, fj fetch.Job) {
	tp := mime.TypeByExtension(path.Ext(fj.Req.URL.Path))
	logx.Printf("mime by extension is %v", tp)
	if tp == "text/html" || tp == "" {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", tp)
	}

}

// showForm either displays a form for requesting an url
// or it returns the URL´s contents.
func showForm(w http.ResponseWriter, r *http.Request) {

	if util.StaticExtension(r) {
		// Handler also serves the root ("host/...").
		// Http pages without proxification
		// cause relative requests to end here.
		// I.e. /css/welt-online.css
		return
	}

	// on live server => always use https
	if cf.ProdHttpsAdvance && r.URL.Scheme != "https" {
		r.URL.Scheme = "https"
		// r.URL.Host = r.Host
		logx.Debugf(r, "redirecting to https: %v", r.URL.String())
		http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
	}

	_ = r.PostFormValue("impossibleKey")
	frm := FormData{}
	dec := formam.NewDecoder(&formam.DecoderOptions{TagName: "formam"})
	err := dec.Decode(r.Form, &frm)
	if err != nil {
		logx.Debugf(r, "Each form field must be accommodated in the struct!\n %v", err)
		return
	}
	if frm.Url != "" {
		logx.Debugf(r, "\n\nForm was %v\n", util.IndentedDump(frm))
	}
	cf.ProcessOptions = frm.ProcessOptions

	s := struct {
		HtmlTitle    string
		ResponseInfo string

		Protocol    string
		Host        string
		Path        string
		UrlParamKey string

		FormData
	}{
		HtmlTitle:    cf.HtmlTitle,
		ResponseInfo: "",

		Protocol:    "https",
		Host:        r.Host, // not fetch.HostFromReq(r)! but why?
		Path:        cf.ProxifyURI,
		UrlParamKey: cf.UrlParamKey,

		FormData: frm,
	}

	if cf.DevHttpFallback {
		s.Protocol = "http"
	}
	if cf.ProdHttpsAdvance {
		s.Protocol = "https"
	}

	if len(frm.Url) > 0 {

		fj := fetch.Job{URL: frm.Url, AeReq: r}
		fj.Fetch()

		if fj.Err != nil ||
			(fj.Status > 299 && fj.Status < 200) {
			s.ResponseInfo = fmt.Sprintf("%v", fj)
		} else {
			// logx.Debugf(r, fmt.Sprintf("%v",fj))
			writeHeader(w, fj)
			cf.Proxifier.InitRequest(r, fj.Req.URL)
			cf.ProcessOptions = frm.ProcessOptions
			cnt := cf.ProcessLean(fj.Bytes())

			pg := HtmlPage{
				Url:  fj.Req.URL.String(),
				Body: cnt,
			}
			paramstr := r.FormValue("param")
			pg.Param, err = strconv.Atoi(paramstr)
			if err != nil {
				// ignore, assuming 0
				pg.Param = valCustom
			}
			valstr := r.FormValue("val")
			pg.Val, err = strconv.Atoi(valstr)
			if err != nil {
				// ignore, assuming 0
				pg.Val = valCustom
			}
			pg.Put(r)

			fmt.Fprintf(w, "%s \n\n", cnt)
			return
		}

		// clean.Beautify = true // "<a> Linktext without trailing space"

	}

	p := logx.PathToSourceFile()
	p = path.Join(p, "templates/*")
	if logx.IsAppengine() {
		workDir, err := os.Getwd()
		util.CheckErr(err)
		p = path.Join(workDir, "templates", "*") // app engine:
	}
	t, err := template.ParseGlob(p)
	util.CheckErr(err)
	err = t.ExecuteTemplate(w, "layout.html", s)
	util.CheckErr(err)

}

func formRedirector(w http.ResponseWriter, r *http.Request) {

	var msg, rURL string

	rURL = r.FormValue(cf.Proxifier.FormRedirectKey)
	if rURL == "" {
		msg := fmt.Sprintf("no url in redirect key %v", cf.Proxifier.FormRedirectKey)
		logx.Debugf(r, msg)
		w.Write([]byte(msg))
		return
	}

	rURL = fmt.Sprintf("%v?1=2&", rURL)
	for key, vals := range r.Form {
		if key == cf.Proxifier.FormRedirectKey {
			continue
		}
		val := vals[0]
		if cf.DevHttpFallback {
			val = strings.Replace(val, " ", "%20", -1)
		}
		rURL = fmt.Sprintf("%v&%v=%v", rURL, key, val)
	}
	logx.Debugf(r, "form redirect url: %v", rURL)

	fj := fetch.Job{URL: rURL, AeReq: r}
	fj.Fetch()
	util.CheckErr(fj.Err)

	cf.Proxifier.InitRequest(r, fj.Req.URL)

	// Todo:
	// Find a way to transfer the options from showForm here.
	// Either addtional form field to form rewrite.
	// Or session.
	// cf.ProcessOptions = frm.ProcessOptions
	cnt := cf.ProcessLean(fj.Bytes())

	writeHeader(w, fj)
	fmt.Fprintf(w, "%s \n\n", cnt)
	fmt.Fprintf(w, "%s \n\n", msg)

	pg := HtmlPage{
		Param: valCustom,
		Val:   valCustom,
		Url:   fj.Req.URL.String(),
		Body:  cnt,
	}
	pg.Put(r)

}
