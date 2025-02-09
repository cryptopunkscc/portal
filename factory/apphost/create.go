package apphost

import (
	"context"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/request"
	"github.com/cryptopunkscc/portal/runtime/apphost"
)

func Full(ctx context.Context) api.Cached {
	return apphost.Cached(Invoker(ctx, api.DefaultClient))
}

func Invoker(ctx context.Context, client api.Client) api.Client {
	return apphost.Invoker{
		Client: client,
		Invoke: Invoke,
		Ctx:    ctx,
		Log:    plog.Get(ctx).Type(apphost.Invoker{}),
	}
}

var Invoke = request.Open.Start
