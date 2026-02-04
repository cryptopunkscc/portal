package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
)

func (a *Adapter) Nodes() *NodeClient {
	return &NodeClient{*a.Client}
}

type NodeClient struct {
	astrald.Client
}

func (c *NodeClient) AddEndpoint(ctx *astral.Context, id astral.Identity, endpoint string) error {
	return Call(ctx, c.Client, "nodes.add_endpoint", query.Args{"id": id.String(), "endpoint": endpoint})
}
