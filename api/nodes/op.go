package nodes

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Op(client apphost.Client) OpClient { return OpClient{client.Rpc().Request("localnode", "nodes")} }

type OpClient struct{ rpc.Conn }

func (n OpClient) AddEndpoint(id string, endpoint string) (err error) {
	return rpc.Command(n, "add_endpoint", rpc.Opt{"id": id, "endpoint": endpoint})
}
