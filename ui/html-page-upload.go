package ui

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/monoculum/formam"
	"github.com/pbberlin/dom/clean"
	"github.com/zew/logx"
	"github.com/zew/util"
)

type UploadRequest struct {
	Url    string `formam:"urlx" json:"urlx"`
	ValId  int    `formam:"val_id" json:"val_id"`
	Data   string `formam:"data" json:"data"` // cannot be byte :(
	Submit string `formam:"submit2" json:"-"`
}

type UploadResponse struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
	Url string `json:"url"`
	Key string `json:"key"`
}

func (ur *UploadResponse) RespondJson(w http.ResponseWriter) {
	js, err := json.Marshal(ur)
	util.CheckErr(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func uploadReceiver(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	ur := new(UploadResponse)
	defer func() {
		logx.Debugf(r, ur.Msg)
		ur.RespondJson(w)
	}()

	_ = r.PostFormValue("impossibleKey")
	frm := UploadRequest{}
	dec := formam.NewDecoder(&formam.DecoderOptions{TagName: "formam"})
	err := dec.Decode(r.Form, &frm)
	if err != nil {
		ur.Msg = fmt.Sprintf("Each form field must be accommodated in the struct!\n %v", err)
		return
	}
	if frm.Url != "" && len(frm.Data) > 0 {
		ur.Msg = fmt.Sprintf("\n\nForm was %v\n", util.Ellipsoider(util.IndentedDump(frm), 120))
	} else {
		ur.Msg = fmt.Sprintf("frm.Url or frm.Data was empty; -%v-  -%v- ", len(frm.Url), len(frm.Data))
		return
	}

	localCleaner := clean.GetDefaultConfig()
	// We dont need Proxifier.InitRequest(r, url)
	// Since we disable proxification:
	localCleaner.Apply(func(c *clean.Config) {
		c.ProcessOptions.SkipProxify = true
	})
	cnt := localCleaner.ProcessLean([]byte(frm.Data))

	pg := HtmlPage{
		Val:  frm.ValId,
		Url:  frm.Url,
		Body: cnt,
	}
	key, err := pg.Put(r)
	if err != nil {
		ur.Msg = fmt.Sprintf("put error %v", err)
		return
	}

	// logx.Debugf(r, "%v", util.Ellipsoider(string(pg.Body), 50))

	ur.Ok = true
	ur.Msg = "saved successfully"
	ur.Url = fmt.Sprintf("/show-detail?val=%v&key=%v", pg.Val, key.Encode())
	ur.Key = key.Encode()

}
