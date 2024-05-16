package open

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"github.com/cryptopunkscc/go-astral-js/runner/tray"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
	src string,
	attach bool,
) (err error) {
	r := portal.Runner[target.App]{
		Action:   action,
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
		return goja.Run(ctx, bindings, app, prefix...)
	case typ.Is(target.Frontend):
		return wails.Run(bindings, app, prefix...)
	default:
		return fmt.Errorf("invalid app target: %v", app.Path())
	}
}

var Handlers = rpc.Handlers{
	"ping":      func() {},
	"open":      portal.NewCmdOpener(portal.ResolveApps, action).Open,
	"observe":   appstore.Observe,
	"install":   apps.Install,
	"uninstall": apps.Uninstall,
}

const action = "o"
