package apphost

import (
	"context"
	"github.com/cryptopunkscc/portal/request/query"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/target"
)

var request = query.Request.Start

func Frontend() target.NewApphost {
	return func(ctx context.Context, portal target.Portal_) target.Apphost {
		adapter := apphost.Adapter(ctx, portal.Manifest().Package)
		return apphost.Invoker(ctx, adapter, request)
	}
}

func Backend() target.NewApphost {
	return func(ctx context.Context, portal target.Portal_) target.Apphost {
		apphost.ConnectionsThreshold = 0
		adapter := apphost.Adapter(ctx, portal.Manifest().Package)
		adapter = apphost.WithTimeout(ctx, adapter, portal)
		return apphost.Invoker(ctx, adapter, request)
	}
}

func Default() target.NewApphost {
	return func(ctx context.Context, portal target.Portal_) target.Apphost {
		return apphost.Adapter(ctx, portal.Manifest().Package)
	}
}
