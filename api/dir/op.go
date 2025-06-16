package dir

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Op(client apphost.Client, target ...string) OpClient {
	return OpClient{client, apphost.Target(target...)}
}

type OpClient struct {
	apphost.Client
	target string
}

func (o OpClient) r() rpc.Conn { return o.Rpc().Request(o.target, "dir") }

func (o OpClient) Resolve(alias string) (out *astral.Identity, err error) {
	s, err := rpc.Query[rpc.Json[*astral.Identity]](o.r(), "resolve", rpc.Opt{
		"name": alias,
		"out":  "json",
	})
	if err != nil {
		return
	}
	out = s.Object
	return
}
