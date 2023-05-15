package main

import (
	"astraljs"
	webview2 "astraljs/webview"
	"github.com/webview/webview"
)

func main() {
	app := astraljs.ResolveWebApp()

	w := webview.New(true)
	defer w.Destroy()

	w.SetSize(800, 600, webview.HintNone)
	w.SetTitle(app.Title)

	// inject apphost js client lib
	w.Init(astraljs.AppHostJsClient())

	// set app source code
	w.SetHtml(app.Source)

	// bind apphost adapter to js env
	webview2.Bind(w, astraljs.NewAppHostFlatAdapter())

	// start js application frontend
	w.Run()
}
