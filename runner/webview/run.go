package webview

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	binding "github.com/cryptopunkscc/go-astral-js/pkg/js/embed/common"
	"github.com/webview/webview"
)

func Run(ctx context.Context, title, src string) {
	w := New(true)
	defer w.Destroy()

	w.SetSize(800, 600, webview.HintNone)
	w.SetTitle(title)

	// inject apphost js client lib
	w.Init(binding.JsString)

	// set app source code
	w.SetHtml(src)

	// bind apphost adapter to js env
	w.BindApphost(apphost.NewAdapter(ctx, nil))

	// start js application frontend
	w.Run()
}
