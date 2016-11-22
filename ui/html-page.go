package ui

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/zew/logx"
	"github.com/zew/util"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type HtmlPage struct {
	Url    string // The Url gets subdomains normalized to dirs; thus its sortable
	Body   []byte `datastore:",noindex"` // []byte to string conversions cause mem copy
	Val    int    // Val is an app specific id, the same Url can be quoted at different positions
	UnixTs int64
	// T      time.Time // We could drop this, if UnixTs
}

const htmlPageKind = "HtmlPage"

func putExample(w http.ResponseWriter, r *http.Request) {
	pg := HtmlPage{
		Val:  rand.Intn(32168),
		Url:  "subdom1.faz.net/aktuell/second.html",
		Body: []byte("some body to love"),
	}
	pg.Put(r)
}

func (pg *HtmlPage) Key() string {
	return fmt.Sprintf("%09v-%v-%v", pg.Val, pg.Url, pg.UnixTs)
}

func (pg *HtmlPage) Put(r *http.Request) (*datastore.Key, error) {
	if !logx.IsAppengine() {
		return nil, fmt.Errorf("Data Store put is only available on app engine")
	}
	ctx := appengine.NewContext(r)
	u, _ := util.UrlParseImproved(pg.Url)
	pg.Url = util.NormalizeSubdomainsToPath(u)

	if pg.UnixTs < 1 {
		pg.UnixTs = (time.Now().Unix() / 600) * 600
	}

	key := datastore.NewKey(ctx, htmlPageKind, pg.Key(), 0, nil)
	keyComplete, err := datastore.Put(ctx, key, pg)
	if err != nil {
		return nil, err
	}
	return keyComplete, nil
}
