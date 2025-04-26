package keys

import (
	"errors"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Client(rpc rpc.Rpc) Conn { return Conn{rpc.Request("localnode", "keys")} }

type Conn struct{ rpc.Conn }

type createKeyArgs struct {
	Alias  string `query:"alias" cli:"alias a"`
	Format string `query:"format" cli:"format f"`
}

func (c Conn) Create(alias string) (*astral.Identity, error) {
	if alias == "" {
		return nil, errors.New("alias is required")
	}
	args := &createKeyArgs{
		Alias:  alias,
		Format: "json",
	}
	return rpc.Query[*astral.Identity](c, "create_key", args)
}
