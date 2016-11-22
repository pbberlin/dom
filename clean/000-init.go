package clean

import (
	"golang.org/x/net/context"

	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/appengine"

	"github.com/zew/logx"
	"github.com/zew/util"
)

type Proxifier struct {
	aeHostName      string // dynamically set if on AE
	aeAppId         string // dynamically set if on AE
	aeDevServerPort string // usage unclear

	proxyHst  string // localhost:xxxx in development; dom-clean.appspot.com in prod
	remoteHst string // host of proxified url, welt.de or faz.net

	ProxifyURI  string `json:"proxify_uri"`   // main uri - relative to proxy host
	UrlParamKey string `json:"url_param_key"` // sadly, hardcoded in formam, urlkey

	FormRedirectKey   string `json:"form_redirect_key"`  // form field name for redirect target
	FormRedirectorURI string `json:"formredirector_uri"` // uri for form redirects, relative to proxy host

	DedupURI string `json:"dedup_uri"`
	DoDedup  bool   `json:"do_dedup"`

	DevHttpFallback  bool `json:"fallback_to_http_dev_srv"`  // fallback to http  on  dev servers
	ProdHttpsAdvance bool `json:"advance_to_https_live_srv"` // advance  to https on live servers
	// AEDev          bool // appengine local dev

}

type ProcessOptions struct {
	FNamer     func() string `json:"-"`
	AddOutline bool          `json:"add_outline"`
	AddID      bool          `json:"add_id_attr"`
	Beautify   bool          `json:"beautify"` // make pretty at the end, removes <a> linktext trailing space

	SkipProxify               bool `formam:"skip_proxify" json:"skip_proxify"`
	SkipPurgeScriptsAndStyles bool `formam:"skip_purge_scripts_and_styles" json:"skip_purge_scripts_and_styles"`
	SkipImage2Link            bool `formam:"skip_img2link" json:"skip_img2link"`
}

type Config struct {
	HtmlTitle string `json:"html_title"`
	ProcessOptions
	Proxifier
}

func (c *Config) setDefaults() {
	c1 := Config{
		HtmlTitle: "Replace with your title",
		Proxifier: Proxifier{
			UrlParamKey:       "urlkey",
			ProxifyURI:        "/prox",
			FormRedirectorURI: "/redir",
			DedupURI:          "/dedup",
			DoDedup:           false,
			DevHttpFallback:   true,
			ProdHttpsAdvance:  false,
		},
		ProcessOptions: ProcessOptions{
			AddOutline:                true,
			AddID:                     true,
			Beautify:                  true,
			SkipProxify:               false,
			SkipPurgeScriptsAndStyles: false,
			SkipImage2Link:            false,
		},
	}
	c = &c1
}

var conf Config

// Get a copy of default config
func GetDefaultConfig() Config {
	return conf
}

// For outside change
func (c *Config) Apply(opts ...func(c *Config)) {
	for _, opt := range opts {
		opt(c)
	}
}

func init() {

	conf.setDefaults()

	// Load from config file
	fileReader := util.LoadConfig("domcleanconfig.json")
	decoder := json.NewDecoder(fileReader)
	err := decoder.Decode(&conf)
	util.CheckErr(err)
	logx.Printf("\n%#s", util.IndentedDump(conf))
}

//
// Setting the proxy host and the remove host
// r         is the request to the proxy
// remoteUrl is the remote url to be proxied
func (p *Proxifier) InitRequest(r *http.Request, remoteUrl *url.URL) {

	p.proxyHst = r.URL.Host // port included!
	if r.URL.Host == "" {
		p.proxyHst = r.Host // port included!
	}

	if false {
		// InstanceID() reads some AE specific environment
		// variable or file system directory.
		// Outside Appengine, it exits with "Metadata fetch failed"
		if appengine.InstanceID() != "" {
			logx.Printf("ae instance Id is %v")
			var ctx context.Context
			func() {
				defer func() {
					rec := recover()
					msg := fmt.Sprintf("appengine panic: %v\n", rec)
					_ = msg
				}()
				ctx = appengine.NewContext(r)
			}()
			if ctx == nil {
				logx.Fatalf("on appengine - but no context")
			}
			p.aeHostName = appengine.DefaultVersionHostname(ctx)
			p.aeAppId = appengine.AppID(ctx)
			if appengine.IsDevAppServer() {
				// Todo: Get the PORT
				p.proxyHst = "localhost" + ":" + p.aeDevServerPort
			} else {
				p.proxyHst = p.aeAppId + ".appspot.com"
			}
		}
	}

	if strings.Index(remoteUrl.Host, ":") > 0 {
		var err error
		p.remoteHst, _, err = net.SplitHostPort(remoteUrl.Host) //just the hostname
		util.CheckErr(err)
	} else {
		p.remoteHst = remoteUrl.Host
	}
}

func (p *Proxifier) ProxyHost() string {
	if p.proxyHst == "" {
		logx.Fatalf("Proxifier proxyHst uninitialized %+v", p)
	}
	return p.proxyHst
}

func (p *Proxifier) RemoteHost() string {
	if p.remoteHst == "" {
		logx.Fatalf("Proxifier remoteHst uninitialized %+v", p)
	}
	return p.remoteHst
}
