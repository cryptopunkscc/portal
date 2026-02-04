package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
)

func (a *Adapter) Keys() *KeysClient {
	return &KeysClient{*a.Client}
}

type KeysClient struct {
	astrald.Client
}

func (c *KeysClient) CreateKey(ctx *astral.Context, alias string) (*astral.Identity, error) {
	return Receive[*astral.Identity](ctx, c.Client, "keys.create_key", query.Args{"alias": alias})
}
