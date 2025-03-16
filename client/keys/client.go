package keys

import (
	"errors"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

func NewClient() Client { return Client{apphost.Default.Rpc().Request("localnode")} }

type Client struct{ rpc.Conn }

type createKeyArgs struct {
	Alias  string `query:"alias" cli:"alias a"`
	Format string `query:"format" cli:"format f"`
}

func (c Client) CreateKey(alias string) (*astral.Identity, error) {
	if alias == "" {
		return nil, errors.New("alias is required")
	}
	args := &createKeyArgs{
		Alias:  alias,
		Format: "json",
	}
	return rpc.Query[*astral.Identity](c, "keys.create_key", args)
}
