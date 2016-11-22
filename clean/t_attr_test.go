package clean

import (
	"fmt"
	"testing"

	"golang.org/x/net/html"
)

func Test_normalizeAttrs(t *testing.T) {

	type tcT struct {
		inp []html.Attribute
		wnt []html.Attribute
	}

	tcs := []tcT{
		tcT{
			[]html.Attribute{
				html.Attribute{Key: "target", Val: "_blank"},
				html.Attribute{Key: "cfrm", Val: "img"},
				html.Attribute{Key: "\n\t  HrEF ", Val: "/dir1"},
			},
			[]html.Attribute{
				html.Attribute{Key: "cfrm", Val: "img"},
				html.Attribute{Key: "href", Val: "/dir1"},
			},
		},
		tcT{
			[]html.Attribute{
				html.Attribute{Key: "href", Val: "/dir1"},
				html.Attribute{Key: "data-href", Val: "/dir2"},
			},
			[]html.Attribute{
				html.Attribute{Key: "href", Val: "/dir1"},
				html.Attribute{Key: "href1", Val: "/dir2"},
			},
		},
		tcT{
			[]html.Attribute{
				html.Attribute{Key: "src", Val: "/dir1/src.img"},
				html.Attribute{Key: "data-src", Val: "/dir2/src.img"},
				html.Attribute{Key: "srcset", Val: "/dir3/src.img"},
				html.Attribute{Key: "style", Val: `input, select{
	margin: 0px 0px;
    padding: 3px 0px;
    font-size: 13px;
    background-image: url("/dir4/src.img");
}`},
			},
			[]html.Attribute{
				html.Attribute{Key: "src", Val: "/dir1/src.img"},
				html.Attribute{Key: "src1", Val: "/dir2/src.img"},
				html.Attribute{Key: "src2", Val: "/dir3/src.img"},
				html.Attribute{Key: "src3", Val: "/dir4/src.img"},
			},
		},
		tcT{
			[]html.Attribute{
				html.Attribute{Key: "type", Val: "input"},
				html.Attribute{Key: "name", Val: "age"},
				html.Attribute{Key: "value", Val: "17"},
			},
			[]html.Attribute{
				html.Attribute{Key: "type", Val: "input"},
				html.Attribute{Key: "name", Val: "age"},
				html.Attribute{Key: "value", Val: "17"},
			},
		},
		tcT{
			[]html.Attribute{
				html.Attribute{Key: "title", Val: "Mr. Brown"},
				html.Attribute{Key: "alt", Val: "Mr. Brown visits Mr. White"},
				html.Attribute{Key: "xxx", Val: "yyy"},
			},
			[]html.Attribute{
				html.Attribute{Key: "title", Val: "Mr. Brown visits Mr. White"},
			},
		},
		tcT{
			[]html.Attribute{
				html.Attribute{Key: "title", Val: "Mr. Brown"},
				html.Attribute{Key: "alt", Val: "Mr. White"},
				html.Attribute{Key: "xxx", Val: "yyy"},
			},
			[]html.Attribute{
				html.Attribute{Key: "title", Val: "Mr. Brown. Mr. White"},
			},
		},
	}

	for i, tc := range tcs {
		inpstr := fmt.Sprintf("%+v", tc.inp)
		wntstr := fmt.Sprintf("%+v", tc.wnt)
		got := normalizeAttrs(tc.inp)
		got = throwAwayAttrs(got)
		gotstr := fmt.Sprintf("%+v", got)
		// t.Logf("%2v: \ninp %20v - \nwnt %20v \ngot %20v\n", i, inpstr, wntstr, gotstr)
		if gotstr != wntstr {
			t.Errorf("%2v: \ninp %20v - \nwnt %20v \ngot %20v\n", i, inpstr, wntstr, gotstr)
		}
	}

}
