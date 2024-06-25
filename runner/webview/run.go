package webview

import (
	"context"
	"github.com/cryptopunkscc/portal/target/apphost"
	"github.com/cryptopunkscc/portal/target/js/embed/common"
	"github.com/webview/webview"
)

func Run(ctx context.Context, title, src string) {
	w := New(true)
	defer w.Destroy()

	w.SetSize(800, 600, webview.HintNone)
	w.SetTitle(title)

	// inject apphost js client lib
	w.Init(common.JsString)

	// set app source code
	w.SetHtml(src)

	// bind apphost adapter to js env
	w.BindApphost(apphost.NewFactory(nil).NewAdapter(ctx, "src"))

	// start js application frontend
	w.Run()
}
