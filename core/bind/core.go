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
	c := DefaultCoreFactory{}.Create(ctx)
	return c, c.Context
}

type Core interface {
	context.Context
	Apphost
	bind.Process
}

type Context struct {
	context.Context
	Adapter
	*Process
}

func (c *Context) GetCtx() context.Context { return c.Context }

type DefaultCoreFactory struct{}

func (DefaultCoreFactory) Create(ctx context.Context) (c *Context) {
	c = &Context{}
	c.Adapter = Adapter{}
	c.Process, ctx = NewProcess(ctx)
	c.Ctx = ctx
	c.Context = ctx
	c.Cached = *apphost.NewCached(apphost.Default)
	c.Adapter.Log = c.Process.log
	return
}

type AutoTokenCoreFactory struct {
	PkgName string
	Tokens  *token.Repository
}

func (f AutoTokenCoreFactory) Create(ctx context.Context) (c Context) {
	t, err := f.Tokens.Get(f.PkgName)
	if err != nil {
		panic(err)
	}

	c.Process, ctx = NewProcess(ctx)
	c.Cached = *apphost.NewCached(apphost.Default)
	c.Ctx = ctx
	c.Context = ctx
	c.Token = t.Token.String()
	c.Endpoint = f.Tokens.Endpoint
	c.Adapter.Log = c.Process.log
	return c
}

func (f AutoTokenCoreFactory) Create2(ctx context.Context, portal target.Portal_) (Core, context.Context) {
	f.PkgName = portal.Manifest().Package
	c := f.Create(ctx)
	return &c, c.Context
}
