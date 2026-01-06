package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/astrald/mod/user"
)

func (a *Adapter) User() *UserClient {
	return &UserClient{*a.Client}
}

type UserClient struct {
	astrald.Client
}

func (c *UserClient) Siblings(ctx *astral.Context) (out <-chan astral.Identity, err error) {
	return GoChan[astral.Identity](ctx, c.Client, "user.list_siblings", query.Args{"zone": astral.ZoneAll})
}

func (c *UserClient) Info(ctx *astral.Context) (out *user.Info, err error) {
	return Receive[*user.Info](ctx, c.Client, "user.info", nil)
}

func (c *UserClient) Claim(ctx *astral.Context, alias string) (out *user.SignedNodeContract, err error) {
	return Receive[*user.SignedNodeContract](ctx, c.Client, "user.claim", query.Args{"target": alias})
}

func (c *UserClient) Create(ctx *astral.Context, alias string) (out *user.CreatedUserInfo, err error) {
	return Receive[*user.CreatedUserInfo](ctx, c.Client, "user.create", query.Args{"alias": alias})
}
