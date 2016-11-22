package clean

var unwantedAttrs = map[string]bool{

	"data": true, // but not data-src, data-href

	"border": false, // check this

	"style": false,
	"class": false,
	// "alt":                 false,
	// "title":               false,

	"align":       false,
	"placeholder": false,

	"target":   false,
	"id":       false,
	"rel":      false,
	"tabindex": false,
	"headline": false,

	"onload":      false,
	"onclick":     false,
	"onmousedown": false,
	"onerror":     false,
	"onsubmit":    false,

	"readonly":       false,
	"accept-charset": false,

	"itemprop":  false,
	"itemtype":  false,
	"itemscope": false,

	"datetime":               false,
	"current-time":           false,
	"fb-iframe-plugin-query": false,
	"fb-xfbml-state":         false,

	"frameborder":       false,
	"async":             false,
	"charset":           false,
	"http-equiv":        false,
	"allowtransparency": false,
	"allowfullscreen":   false,
	"scrolling":         false,
	"ftghandled":        false,
	"ftgrandomid":       false,
	"marginwidth":       false,
	"marginheight":      false,
	"vspace":            false,
	"hspace":            false,
	"seamless":          false,
	"aria-hidden":       false,
	"gapi_processed":    false,
	"property":          false,
	"media":             false,

	"content":  false,
	"language": false,

	"role":  false,
	"sizes": false,
}
