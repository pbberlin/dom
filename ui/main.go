// package ui implements http user interface for dom/clean;
// it is 'plugabble' into other go programmes.
package ui

import (
	"log"
	"net/http"
	"reflect"

	"github.com/pbberlin/dom/clean"
	"github.com/zew/logx"
)

func main() {
	log.Println("Use Test_Build() instead of build")
}

var cf clean.Config = clean.GetDefaultConfig()

// Explicitly callable
// Used by the tests
func ExplicitInit(mux *http.ServeMux) {

	if mux == nil {
		http.HandleFunc("/", showForm)
		http.HandleFunc(cf.ProxifyURI, showForm)
		http.HandleFunc(cf.FormRedirectorURI, formRedirector)
		http.HandleFunc("/put-example", putExample)
		http.HandleFunc("/query-pages", queryPages)
		http.HandleFunc("/show-page", showPage)
		http.HandleFunc("/show-page-by-url", showPageByUrl)
		http.HandleFunc("/delete-page", deletePage)
		http.HandleFunc("/upload-receiver", uploadReceiver)
	} else {
		// mux.HandleFunc("/", showForm)
		mux.HandleFunc(cf.ProxifyURI, showForm)
		mux.HandleFunc(cf.FormRedirectorURI, formRedirector)
		mux.HandleFunc("/put-example", putExample)
		mux.HandleFunc("/query-pages", queryPages)
		mux.HandleFunc("/show-page", showPage)
		mux.HandleFunc("/delete-page", deletePage)
		mux.HandleFunc("/upload-receiver", uploadReceiver)
	}

	fd := new(FormData)
	val := reflect.ValueOf(fd).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		_ = valueField
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if typeField.Name == "Url" && tag.Get("formam") != cf.UrlParamKey {
			logx.Fatalf("Formam tag for Url must match cf.UrlParamKey.\nInstead %v vs %v",
				tag.Get("formam"), cf.UrlParamKey,
			)
		}
		// logx.Printf("Field Name: %q,\t Field Value: -%v-,\t Tag Value: %q\n",
		// 	typeField.Name, valueField.Interface(), tag.Get("formam"),
		// )
	}
}

func init() {
	// ExplicitInit()
}
