package keys

import (
	"errors"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Client(rpc rpc.Rpc) Conn { return Conn{rpc.Request("localnode", "keys")} }

type Conn struct{ rpc.Conn }

type createKeyArgs struct {
	Alias string `query:"alias" cli:"alias a"`
	Out   string `query:"out" cli:"out o"`
}

func (c Conn) Create(alias string) (id *astral.Identity, err error) {
	if alias == "" {
		return nil, errors.New("alias is required")
	}
	args := &createKeyArgs{
		Alias: alias,
		Out:   "json",
	}
	o, err := rpc.Query[rpc.JsonObject[*astral.Identity]](c, "create_key", args)
	if err != nil {
		return
	}
	id = o.Object
	return
}
