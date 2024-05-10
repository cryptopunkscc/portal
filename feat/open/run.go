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
	r := portal.Runner[target.App]{
		Action:   "open",
		Port:     "portal",
		New:      bindings,
		Tray:     tray.Run,
		Serve:    serve.Run,
		Resolve:  portal.ResolveApps,
		Attach:   Attach,
		Handlers: Handlers,
	}
	return r.Run(ctx, src, attach)
}

func Attach(
	ctx context.Context,
	bindings runtime.New,
	app target.App,
	prefix ...string,
) (err error) {
	typ := app.Type()
	switch {
	case typ.Is(target.Backend):
		if err = goja.NewBackend(bindings(target.Backend, prefix...)).RunFs(app.Files()); err != nil {
			return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
		}
		<-ctx.Done()
	case typ.Is(target.Frontend):
		opt := wails.AppOptions(bindings(target.Frontend, prefix...))
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
