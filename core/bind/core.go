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

type DefaultCoreFactory struct {
	Adapter *apphost.Adapter
}

func (f DefaultCoreFactory) Create(ctx context.Context) (c *Context) {
	if f.Adapter == nil {
		f.Adapter = apphost.Default
	}
	c = &Context{}
	c.Process, c.Context = NewProcess(ctx)
	c.Ctx = c.Context
	c.Cached = *apphost.NewCached(f.Adapter)
	c.Adapter.Log = c.Process.log
	return
}

type AutoTokenCoreFactory struct {
	PkgName string
	Tokens  *token.Repository
}

func (f AutoTokenCoreFactory) Create(ctx context.Context) (c *Context) {
	c = DefaultCoreFactory{}.Create(ctx)

	t, err := f.Tokens.Get(f.PkgName)
	if err != nil {
		panic(err)
	}

	c.Token = t.Token.String()
	c.Endpoint = f.Tokens.Endpoint
	return c
}

func (f AutoTokenCoreFactory) Create2(ctx context.Context, portal target.Portal_) (Core, context.Context) {
	f.PkgName = portal.Manifest().Package
	c := f.Create(ctx)
	return c, c.Context
}
