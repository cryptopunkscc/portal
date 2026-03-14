package bind

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/client"
)

type Core struct {
	context.Context
	Astrald
	*Process
}

func (c *Core) GetCtx() context.Context { return c.Context }

type DefaultCoreFactory struct {
	Astrald *client.Astrald
}

func (f DefaultCoreFactory) Create(ctx context.Context) (c *Core) {
	if f.Astrald == nil {
		f.Astrald = client.Default
	}
	c = &Core{}
	c.Process, c.Context = NewProcess(ctx)
	c.Astrald.Log = c.Process.log
	c.Init()
	return
}
