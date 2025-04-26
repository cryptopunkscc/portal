package nodes

import "github.com/cryptopunkscc/portal/pkg/rpc"

func Client(rpc rpc.Rpc) Conn { return Conn{rpc.Request("localnode", "nodes")} }

type Conn struct{ rpc.Conn }

func (n Conn) AddEndpoint(id string, endpoint string) (err error) {
	return rpc.Command(n, "add_endpoint", rpc.Opt{"id": id, "endpoint": endpoint})
}
