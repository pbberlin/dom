package clean

import (
	"fmt"
	"strings"

	"github.com/pbberlin/dom/node"
	"github.com/zew/util"

	"golang.org/x/net/html"
)

// Take a relative URI from href, src or action
// and prefix it with the (remote) hostname.
// Special case: Links to images do not need https.
// Use forceHttp to force http
func (p *Proxifier) absolutize(host, uri string, forceHttp bool) string {
	if strings.HasPrefix(uri, "//ssl.") || !strings.HasPrefix(uri, "/") {
		return uri
	}
	proto := "https"
	if p.DevHttpFallback {
		proto = "http"
	} else if p.ProdHttpsAdvance {
		proto = "https"
	}
	if forceHttp {
		proto = "http"
	}
	uri = fmt.Sprintf("%v://%v%v", proto, host, uri)
	return uri
}

// Rewrite href, src and action
// for some node types anchor, image, form
func (p *Proxifier) attrsAbsoluteAndProxified(attrs []html.Attribute) []html.Attribute {

	rew := make([]html.Attribute, 0, len(attrs))
	for i := 0; i < len(attrs); i++ {
		attr := attrs[i]
		attr.Val = strings.TrimSpace(attr.Val)
		switch attr.Key {
		case "href", "src", "data-src", "srcset":
			if util.IsSpacey(attr.Val) {
				continue // throw away
			}
			startsWithSlash := strings.HasPrefix(attr.Val, "/")
			startsWithRune := strings.HasPrefix(attr.Val, "#")
			isOrWasImg := (node.AttrX(attrs, "cfrom") == "img") || (attr.Key == "src") || (attr.Key == "data-src")
			isOrWasImg = isOrWasImg || util.ImageExtension(attr.Val)
			// Make absolute
			if startsWithSlash && !startsWithRune {
				attr.Val = p.absolutize(p.RemoteHost(), attr.Val, isOrWasImg)
			}
			// Proxify
			if !isOrWasImg {
				if p.DoDedup {
					attr.Val = fmt.Sprintf("%v?%v=%v", p.DedupURI, p.UrlParamKey, attr.Val)
				} else {
					attr.Val = fmt.Sprintf("%v?%v=%v", p.ProxifyURI, p.UrlParamKey, attr.Val)
				}
			}
		case "action":
			attr.Val = p.absolutize(p.ProxyHost(), p.FormRedirectorURI, false)
		case "method":
			attr.Val = "post"
		}
		rew = append(rew, attr)
	}

	// We instrumented all forms with a field "redirect-to"
	// Now we have to make the value of this field absolute
	isRedir := node.AttrX(attrs, p.FormRedirectKey)
	if isRedir != "" {
		for i := 0; i < len(rew); i++ {
			if rew[i].Key == "value" {
				// rew[i].Key = p.absolutize(p.RemoteHost(), rew[i].Key, false)
			}
		}
	}

	rew = append(rew, html.Attribute{Key: "attrs", Val: "rewritten"})
	return rew
}
