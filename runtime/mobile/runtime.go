package runtime

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	apphostRuntime "github.com/cryptopunkscc/portal/runtime/apphost"
	bindRuntime "github.com/cryptopunkscc/portal/runtime/bind"
)

type Mobile struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Serve  mobile.Serve
	Client apphost.Client
	Find   target.Find[target.App_]
}

func (m *Mobile) Stop()                   { m.Cancel() }
func (m *Mobile) Start()                  { m.Serve() }
func (m *Mobile) Apphost() mobile.Apphost { return Apphost(m.Client) }
func (m *Mobile) App(path string) mobile.App {
	apps, err := m.Find(m.Ctx, path)
	if err != nil {
		plog.Get(m.Ctx).Type(m).E().Println(err)
		return nil
	}
	if len(apps) == 0 {
		plog.Get(m.Ctx).Type(m).E().Println(target.ErrNotFound)
		return nil
	}
	app := apps[0]
	return App{
		source: app,
		runtime: &bind.Module{
			Sys: bindRuntime.Sys(m.Ctx),
			Apphost: bindRuntime.Adapter(
				m.Ctx,
				apphostRuntime.Default(),
				app.Manifest().Package,
			),
		},
	}
}
