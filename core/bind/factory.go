package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/target"
	core "github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/request"
)

func DefaultCore() NewCore  { return api(frontendApphost(cachedInvoker)) }
func FrontendCore() NewCore { return api(frontendApphost(cachedInvoker)) }
func BackendCore() NewCore  { return api(backendApphost(cachedInvoker)) }

func api(newApphost newApphost) NewCore {
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

func cachedInvoker(ctx context.Context) apphost.Cached {
	return core.Cached(core.Invoker{
		Client: core.Default,
		Invoke: request.Open,
		Ctx:    ctx,
		Log:    plog.Get(ctx).Type(core.Invoker{}),
	})
}

type newApphost func(ctx context.Context, portal target.Portal_) Apphost

func frontendApphost(create func(ctx context.Context) apphost.Cached) newApphost {
	return func(ctx context.Context, portal target.Portal_) Apphost {
		return Adapter(ctx, create(ctx), portal.Manifest().Package)
	}
}

func backendApphost(create func(ctx context.Context) apphost.Cached) newApphost {
	return func(ctx context.Context, portal target.Portal_) Apphost {
		cached := create(ctx)
		core.ConnectionsThreshold = 0
		core.Timeout(ctx, cached, portal)
		return Adapter(ctx, cached, portal.Manifest().Package)
	}
}
