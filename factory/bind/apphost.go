package bind

import (
	"context"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	apphost2 "github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/bind"
)

type NewApphost func(ctx context.Context, portal target.Portal_) bind.Apphost

var DefaultApphost = FrontendApphost

func FrontendApphost(create func(ctx context.Context) api.Cached) NewApphost {
	return func(ctx context.Context, portal target.Portal_) bind.Apphost {
		return bind.Adapter(ctx, create(ctx), portal.Manifest().Package)
	}
}

func BackendApphost(create func(ctx context.Context) api.Cached) NewApphost {
	return func(ctx context.Context, portal target.Portal_) bind.Apphost {
		apphost := create(ctx)
		apphost2.ConnectionsThreshold = 0
		apphost2.Timeout(ctx, apphost, portal)
		return bind.Adapter(ctx, apphost, portal.Manifest().Package)
	}
}
