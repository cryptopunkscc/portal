package apphost

import (
	"errors"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func (a *Adapter) Key() Key { return Key{a.Rpc().Request("localnode")} }

type Key struct{ rpc.Conn }

type createKeyArgs struct {
	Alias  string `query:"alias" cli:"alias a"`
	Format string `query:"format" cli:"format f"`
}

func (c Key) Create(alias string) (*astral.Identity, error) {
	if alias == "" {
		return nil, errors.New("alias is required")
	}
	args := &createKeyArgs{
		Alias:  alias,
		Format: "json",
	}
	return rpc.Query[*astral.Identity](c, "keys.create_key", args)
}
