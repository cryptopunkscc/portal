package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/astral"
	. "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/request"
	"github.com/cryptopunkscc/portal/runtime/apphost"
)

var Basic = apphost.Adapter(astral.Client)

func Default() Cached {
	return apphost.Cached(apphost.Adapter(astral.Client))
}

func Full(ctx context.Context) Cached {
	return apphost.Cached(Invoker(ctx, apphost.Adapter(astral.Client)))
}

var Invoke = request.Open.Start

func Invoker(ctx context.Context, client Client) Client {
	return apphost.Invoker{
		Client: client,
		Invoke: Invoke,
		Log:    plog.Get(ctx).Type(apphost.Invoker{}),
		Ctx:    ctx,
	}
}
