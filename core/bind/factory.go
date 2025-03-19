package bind

import (
	"context"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/request"
)

func DefaultRuntime() NewRuntime  { return newRuntime(frontendApphost(cachedInvoker)) }
func FrontendRuntime() NewRuntime { return newRuntime(frontendApphost(cachedInvoker)) }
func BackendRuntime() NewRuntime  { return newRuntime(backendApphost(cachedInvoker)) }

func newRuntime(newApphost newApphost) NewRuntime {
	return func(ctx context.Context, portal target.Portal_) (Runtime, context.Context) {
		m := Module{}
		m.Sys, ctx = Sys(ctx)
		m.Apphost = newApphost(ctx, portal)
		return m, ctx
	}
}

func cachedInvoker(ctx context.Context) api.Cached {
	return apphost.Cached(apphost.Invoker{
		Client: apphost.Default,
		Invoke: request.Open,
		Ctx:    ctx,
		Log:    plog.Get(ctx).Type(apphost.Invoker{}),
	})
}

type newApphost func(ctx context.Context, portal target.Portal_) Apphost

func frontendApphost(create func(ctx context.Context) api.Cached) newApphost {
	return func(ctx context.Context, portal target.Portal_) Apphost {
		return Adapter(ctx, create(ctx), portal.Manifest().Package)
	}
}

func backendApphost(create func(ctx context.Context) api.Cached) newApphost {
	return func(ctx context.Context, portal target.Portal_) Apphost {
		cached := create(ctx)
		apphost.ConnectionsThreshold = 0
		apphost.Timeout(ctx, cached, portal)
		return Adapter(ctx, cached, portal.Manifest().Package)
	}
}
