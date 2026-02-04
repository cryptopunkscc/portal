package bind

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/apphost"
)

type Core struct {
	context.Context
	Adapter
	*Process
}

func (c *Core) GetCtx() context.Context { return c.Context }

type DefaultCoreFactory struct {
	Adapter *apphost.Adapter
}

func (f DefaultCoreFactory) Create(ctx context.Context) (c *Core) {
	if f.Adapter == nil {
		f.Adapter = apphost.Default
	}
	c = &Core{}
	c.Process, c.Context = NewProcess(ctx)
	c.Adapter.Log = c.Process.log
	c.Init()
	return
}
