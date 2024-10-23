package bind

import (
	"context"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/request/query"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/target"
)

var Request = query.Request.Start

type NewApphost func(ctx context.Context, portal target.Portal_) bind.Apphost

func FrontendApphost(apphost api.Cached) NewApphost {
	return func(ctx context.Context, portal target.Portal_) bind.Apphost {
		adapter := bind.Adapter(ctx, apphost, portal.Manifest().Package)
		return bind.ApphostInvoker(ctx, adapter, Request)
	}
}

func BackendApphost(apphost api.Cached) NewApphost {
	return func(ctx context.Context, portal target.Portal_) bind.Apphost {
		bind.ConnectionsThreshold = 0
		adapter := bind.Adapter(ctx, apphost, portal.Manifest().Package)
		adapter = bind.ApphostTimeout(ctx, adapter, portal)
		return bind.ApphostInvoker(ctx, adapter, Request)
	}
}

func DefaultApphost(apphost api.Cached) NewApphost {
	return func(ctx context.Context, portal target.Portal_) bind.Apphost {
		return bind.Adapter(ctx, apphost, portal.Manifest().Package)
	}
}
