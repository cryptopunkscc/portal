package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/request/query"
	runtime "github.com/cryptopunkscc/portal/runtime/apphost"
)

func Full(ctx context.Context) apphost.Cached { return runtime.Cached(Invoker(ctx)) }

func Cached() apphost.Cached { return runtime.Cached(Client()) }

func Invoker(ctx context.Context) apphost.Client {
	i := runtime.Invoker{}
	i.Client = Client()
	i.Invoke = Invoke
	i.Log = plog.Get(ctx).Type(i).Set(&ctx)
	i.Ctx = ctx
	return i
}

func Client() apphost.Client { return runtime.Adapter(astral.Client) }

var Invoke = query.Open.Start
