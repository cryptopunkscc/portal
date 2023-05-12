package main

import (
	"astral-js"
	webview2 "astral-js/webview"
	"github.com/webview/webview"
)

func main() {
	app := astral_js.ResolveWebApp()

	w := webview.New(true)
	defer w.Destroy()

	w.SetSize(800, 600, webview.HintNone)
	w.SetTitle(app.Title)

	// inject apphost js client lib
	w.Init(astral_js.AppHostJsClient())

	// set app source code
	w.SetHtml(app.Source)

	// bind apphost adapter to js env
	webview2.Bind(w, astral_js.NewAppHostFlatAdapter())

	// start js application frontend
	w.Run()
}
