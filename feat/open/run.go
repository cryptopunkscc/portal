package open

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/feat/tray"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
	src string,
	attach bool,
) (err error) {
	return portal.Runner[target.App]{
		Cmd:      "open",
		New:      bindings,
		Tray:     tray.Run,
		Serve:    serve.Run,
		Resolve:  portal.ResolveApps,
		Attach:   Attach,
		Handlers: Handlers,
	}.Run(ctx, src, attach)
}

func Attach(
	ctx context.Context,
	bindings runtime.New,
	app target.App,
) (err error) {
	switch app.Type() {

	case target.Backend:
		if err = goja.NewBackend(bindings(app.Type())).RunFs(app.Files()); err != nil {
			return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
		}
		<-ctx.Done()

	case target.Frontend:
		opt := wails.AppOptions(bindings(app.Type()))
		if err = wails.Run(app, opt); err != nil {
			return fmt.Errorf("dev.Run: %v", err)
		}

	default:
		return fmt.Errorf("invalid target: %v", app.Path())
	}
	return
}

var Handlers = rpc.Handlers{
	"ping":      func() {},
	"open":      portal.NewCmdOpener(portal.ResolveApps, "open").Open,
	"observe":   appstore.Observe,
	"install":   appstore.Install,
	"uninstall": appstore.Uninstall,
}
