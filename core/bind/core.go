package bind

import (
	"context"

	lib "github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/core/token"
)

type NewCore func(ctx context.Context, portal target.Portal_) (Core, context.Context)

func CreateCore(ctx context.Context, portal target.Portal_) (Core, context.Context) {
	return DefaultCoreFactory{}.Create(ctx)
}

type Core interface {
	Apphost
	bind.Process
}

type core struct {
	Apphost
	bind.Process
}

type DefaultCoreFactory struct{}

func (DefaultCoreFactory) Create(ctx context.Context) (Core, context.Context) {
	i := &apphost.Invoker{Ctx: ctx}
	c := core{}
	c.Process, ctx = Process(ctx)
	c.Apphost = Adapter(ctx, apphost.NewCached(i))
	return c, ctx
}

type LibCoreFactory struct{ Client *lib.Client }

func (f LibCoreFactory) Create(ctx context.Context) (Core, context.Context) {
	i := &apphost.Invoker{Ctx: ctx}
	i.Client.Endpoint = f.Client.Endpoint
	i.Client.AuthToken = f.Client.AuthToken
	c := core{}
	c.Process, ctx = Process(ctx)
	c.Apphost = Adapter(ctx, apphost.NewCached(i))
	return c, ctx
}

type AutoTokenCoreFactory struct {
	PkgName string
	Tokens  *token.Repository
}

func (f AutoTokenCoreFactory) Create(ctx context.Context) (Core, context.Context) {
	i := &apphost.Invoker{Ctx: ctx}
	i.Client.Endpoint = f.Tokens.Endpoint
	t, err := f.Tokens.Get(f.PkgName)
	if err != nil {
		return nil, nil
	}
	i.Client.AuthToken = t.Token.String()
	c := core{}
	c.Process, ctx = Process(ctx)
	c.Apphost = Adapter(ctx, apphost.NewCached(i))
	return c, ctx
}

func (f AutoTokenCoreFactory) Create2(ctx context.Context, portal target.Portal_) (Core, context.Context) {
	f.PkgName = portal.Manifest().Package
	return f.Create(ctx)
}
