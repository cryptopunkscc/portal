package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/target"
	core "github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/request"
)

func NewDefaultCoreFunc() NewCore  { return CoreFactory{}.NewDefaultFunc() }
func NewFrontendCoreFunc() NewCore { return CoreFactory{}.NewFrontendFunc() }
func NewBackendCoreFunc() NewCore  { return CoreFactory{}.NewBackendFunc() }

type CoreFactory struct{ token.Repository }

func (f CoreFactory) NewDefaultFunc() NewCore  { return f.api(f.frontendApphost(f.cachedInvoker)) }
func (f CoreFactory) NewFrontendFunc() NewCore { return f.api(f.frontendApphost(f.cachedInvoker)) }
func (f CoreFactory) NewBackendFunc() NewCore  { return f.api(f.frontendApphost(f.cachedInvoker)) }

func (f CoreFactory) api(newApphost newApphost) NewCore {
	return func(ctx context.Context, portal target.Portal_) (Core, context.Context) {
		m := struct {
			Apphost
			bind.Process
		}{}
		m.Process, ctx = Process(ctx)
		m.Apphost = newApphost(ctx, portal)
		return m, ctx
	}
}

func (f CoreFactory) cachedInvoker(ctx context.Context, portal target.Portal_) apphost.Cached {
	a := &core.Adapter{}
	if f.Adapter != nil {
		a.Endpoint = f.Adapter.Endpoint
	}
	if t, err := f.Repository.Get(portal.Manifest().Package); err == nil {
		a.AuthToken = string(t.Token)
	}
	return core.Cached(core.Invoker{
		Client: a,
		Invoke: request.Open,
		Ctx:    ctx,
		Log:    plog.Get(ctx).Type(core.Invoker{}),
	})
}

type newApphost func(ctx context.Context, portal target.Portal_) Apphost

func (f CoreFactory) frontendApphost(
	create func(ctx context.Context, portal target.Portal_) apphost.Cached,
) newApphost {
	return func(ctx context.Context, portal target.Portal_) Apphost {
		return Adapter(ctx, create(ctx, portal), portal.Manifest().Package)
	}
}

func (f CoreFactory) backendApphost(
	create func(ctx context.Context, portal target.Portal_) apphost.Cached,
) newApphost {
	return func(ctx context.Context, portal target.Portal_) Apphost {
		cached := create(ctx, portal)
		core.ConnectionsThreshold = 0
		core.Timeout(ctx, cached, portal)
		return Adapter(ctx, cached, portal.Manifest().Package)
	}
}
