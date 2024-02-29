package main

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	webview2 "github.com/cryptopunkscc/go-astral-js/pkg/runner/frontend/webview"
	"github.com/webview/webview"
	"log"
	"os"
	"path"
)

func main() {
	file := os.Args[1]

	srcBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	title := path.Base(file)
	src := string(srcBytes)

	w := webview.New(true)
	defer w.Destroy()

	w.SetSize(800, 600, webview.HintNone)
	w.SetTitle(title)

	// inject apphost js client lib
	w.Init(apphost.JsBaseString())

	// set app source code
	w.SetHtml(src)

	// bind apphost adapter to js env
	webview2.Bind(w, apphost.NewFlatAdapter())

	// start js application frontend
	w.Run()
}
