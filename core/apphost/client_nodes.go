package apphost

import "github.com/cryptopunkscc/portal/pkg/rpc"

func (a *Adapter) Nodes() Nodes { return Nodes{a.Rpc().Request("localnode", "nodes")} }

type Nodes struct{ rpc.Conn }

func (n Nodes) AddEndpoint(id string, endpoint string) (err error) {
	return rpc.Command(n, "add_endpoint", rpc.Opt{"id": id, "endpoint": endpoint})
}
