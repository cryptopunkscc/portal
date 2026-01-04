package bind

import (
	"context"

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

type Core2 struct {
	*Adapter
	*Process
}

type DefaultCoreFactory struct{}

func (DefaultCoreFactory) Create(ctx context.Context) (Core2, context.Context) {
	i := &apphost.Invoker{Ctx: ctx}
	c := Core2{}
	c.Process, ctx = NewProcess(ctx)
	c.Adapter = NewAdapter(ctx, apphost.NewCached(i))
	return c, ctx
}

type AutoTokenCoreFactory struct {
	PkgName string
	Tokens  *token.Repository
}

func (f AutoTokenCoreFactory) Create(ctx context.Context) (Core2, context.Context) {
	i := &apphost.Invoker{Ctx: ctx}
	i.Endpoint = f.Tokens.Endpoint
	t, err := f.Tokens.Get(f.PkgName)
	if err != nil {
		return Core2{}, nil
	}
	i.Token = t.Token.String()
	c := Core2{}
	c.Process, ctx = NewProcess(ctx)
	c.Adapter = NewAdapter(ctx, apphost.NewCached(i))
	return c, ctx
}

func (f AutoTokenCoreFactory) Create2(ctx context.Context, portal target.Portal_) (Core, context.Context) {
	f.PkgName = portal.Manifest().Package
	return f.Create(ctx)
}
