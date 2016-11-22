// +build parsing

// go test -tags=parsing
package clean

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"appengine/aetest"

	"github.com/pbberlin/tools/net/http/fileserver"
	"github.com/pbberlin/tools/net/http/repo"
	"github.com/pbberlin/tools/net/http/routes"
	"github.com/pbberlin/tools/sort/sortmap"
	"github.com/zew/exceldb/fetch"
	"github.com/zew/logx"
)

const numTotal = 3 // comparable html docs
const stageMax = 3 // weedstages

const cTestHostOwn = "localhost:63222"

var hostWithPref = routes.AppHost() + repo.UriMountNameY

func prepare(t *testing.T) aetest.Context {

	c, err := aetest.NewContext(nil)
	if err != nil {
		logx.Printf("%v", err)
		t.Fatal(err)
	}

	serveFile := func(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
		fs1 := repo.GetFS(c)
		fileserver.FsiFileServer(w, r, fileserver.Options{FS: fs1, Prefix: repo.UriMountNameY})
	}
	http.HandleFunc(repo.UriMountNameY, serveFile)

	go func() {
		log.Fatal(
			http.ListenAndServe(cTestHostOwn, nil),
		)
	}()

	return c

}

func Test1(t *testing.T) {

	// c := prepare(t)
	// defer c.Close()

	logx.Printf("waiting for webserver")
	time.Sleep(2 * time.Millisecond)

	remoteHostname := "www.welt.de"

	dirs1, _, msg, err := fileserver.GetDirContents(hostWithPref, remoteHostname)
	if err != nil {
		logx.Printf("%v", err)
		logx.Printf("%s", msg)
	}

	logx.Printf("dirs1")
	for _, v := range dirs1 {
		logx.Printf("    %v", v)
	}

	least3Files := []string{}
	for _, v1 := range dirs1 {

		dirs2, fils2, msg, err := fileserver.GetDirContents(hostWithPref, path.Join(remoteHostname, v1))
		_ = dirs2
		if err != nil {
			logx.Printf("%v", err)
			logx.Printf("%s", msg)
		}
		// logx.Printf("  dirs2 %v", util.IndentedDump(dirs2))
		// logx.Printf("  fils2 %v", util.IndentedDump(fils2))

		if len(fils2) > numTotal-1 {
			for i2, v2 := range fils2 {
				least3Files = append(least3Files, path.Join(remoteHostname, v1, v2))
				if i2 == numTotal-1 {
					break
				}
			}
			break
		}
	}

	if len(least3Files) < numTotal {
		logx.Printf("not enough files in rss fetcher cache")
		return
	}

	logx.Printf("fils2")
	for _, v := range least3Files {
		logx.Printf("    %v", v)
	}

	logdir := prepareLogDir()

	iter := make([]int, numTotal)

	for i, _ := range iter {

		surl := spf("%v/%v", hostWithPref, least3Files[i])

		fNamer := FileNamer(logdir, i)
		fnKey := fNamer() // first call yields key
		_ = fnKey

		fj := fetch.Job{URL: surl, AeReq: nil}
		fj.Fetch()
		if fj.Err != nil {
			logx.Printf("%v", fj.Err)
			return
		}
		resBytes := fj.Bytes()
		effUrl := fj.Req.URL

		logx.Printf("fetched %4.1fkB from %v", float64(len(resBytes))/1024, util.ToLenR(effUrl.String(), 60))
		opts := Config{}
		opts.FNamer = fNamer
		opts.RemoteHost = remoteHostname
		doc, err := DomClean(resBytes, opts)
		logx.Printf("%v", err)
		_ = doc

	}

	// statistics on elements and attributes
	sorted1 := sortmap.SortMapByCount(attrDistinct)
	sorted1.Print(6)
	fmt.Println()
	sorted2 := sortmap.SortMapByCount(nodeDistinct)
	sorted2.Print(6)

	pf("correct finish\n")

}

func prepareLogDir() string {

	logdir := "outp"
	logx.Printf("logdir is %v ", logdir)

	// sweep previous
	rmPath := spf("./%v/", logdir)
	err := os.RemoveAll(rmPath)
	if err != nil {
		logx.Printf("%v", err)
		os.Exit(1)
	}
	logx.Printf("removed %q", rmPath)

	// create anew
	err = os.Mkdir(logdir, 0755)
	if err != nil && !os.IsExist(err) {
		logx.Printf("%v", err)
		os.Exit(1)
	}

	return logdir

}
