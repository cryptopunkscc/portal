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
		adapter := apphost.NewAdapter(ctx, portal.Manifest().Package)
		return apphost.NewInvoker(ctx, adapter, request)
	}
}

func Backend() target.NewApphost {
	return func(ctx context.Context, portal target.Portal_) target.Apphost {
		apphost.ConnectionsThreshold = 0
		adapter := apphost.NewAdapter(ctx, portal.Manifest().Package)
		adapter = apphost.WithTimeout(ctx, adapter, portal)
		return apphost.NewInvoker(ctx, adapter, request)
	}
}

func Default() target.NewApphost {
	return func(ctx context.Context, portal target.Portal_) target.Apphost {
		return apphost.NewAdapter(ctx, portal.Manifest().Package)
	}
}
