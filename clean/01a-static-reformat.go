package clean

import (
	"bytes"
	"strings"

	"github.com/zew/util"
)

var insertNewlines = strings.NewReplacer(
	"<head", "\n<head",
	"</head>", "</head>\n",

	"</script>", "</script>\n",
	"<script", "\n<script",

	"</noscript>", "</noscript>\n",
	"<noscript", "\n<noscript",

	"</style>", "</style>\n",
	"<style", "\n<style",

	"</link>", "</link>\n",
	"<link", "\n<link",

	"-->", "-->\n",
	"<!--", "\n<!--",

	"<meta", "\n<meta",
	"</div>", "</div>\n",
)

func globalByteFixes(b []byte) []byte {
	// <!--(.*?)-->
	b = bytes.Replace(b, []byte("<!--<![endif]-->"), []byte("<![endif]-->"), -1)
	return b
}

func primitiveReformat(s string) string {
	s = insertNewlines.Replace(s)
	return util.UndoubleNewlines(s)
}
