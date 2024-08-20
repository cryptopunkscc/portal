package runtime

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/mobile"
	apphostRuntime "github.com/cryptopunkscc/portal/runtime/apphost"
	bindRuntime "github.com/cryptopunkscc/portal/runtime/bind"
)

type Mobile struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Serve  mobile.Serve
	Client apphost.Client
}

func (m *Mobile) Stop()                   { m.Cancel() }
func (m *Mobile) Start()                  { m.Serve() }
func (m *Mobile) Apphost() mobile.Apphost { return Apphost(m.Client) }
func (m *Mobile) Bindings(pkg string) bind.Runtime {
	return &bind.Module{
		Sys:     bindRuntime.Sys(m.Ctx),
		Apphost: bindRuntime.Adapter(m.Ctx, apphostRuntime.Default(), pkg),
	}
}
