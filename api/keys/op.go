package keys

import (
	"errors"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Op(client apphost.Client) OpClient { return OpClient{client.Rpc().Request("localnode", "keys")} }

type OpClient struct{ rpc.Conn }

func (c OpClient) CreateKey(alias string) (id *astral.Identity, err error) {
	if alias == "" {
		return nil, errors.New("alias is required")
	}
	o, err := rpc.Query[rpc.Json[*astral.Identity]](c, "create_key", rpc.Opt{"out": "json", "alias": alias})
	if err != nil {
		return
	}
	id = o.Object
	return
}
