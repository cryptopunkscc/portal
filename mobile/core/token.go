package core

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/apphost"
	bind2 "github.com/cryptopunkscc/portal/pkg/bind/src"
)

type AutoTokenCoreFactory struct {
	PkgName string
	Tokens  *apphost.Tokens
}

func (f AutoTokenCoreFactory) Create(ctx context.Context) (c *bind2.Core) {
	c = bind2.DefaultCoreFactory{}.Create(ctx)

	t, err := f.Tokens.Get(f.PkgName)
	if err != nil {
		panic(err)
	}

	c.Token = t.Token.String()
	c.Endpoint = f.Tokens.Adapter.Endpoint
	return c
}
