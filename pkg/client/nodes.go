package client

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
)

type Nodes struct {
	astrald.Client
}

func (c *Nodes) AddEndpoint(ctx *astral.Context, id astral.Identity, endpoint string) error {
	return Call(ctx, c.Client, "nodes.add_endpoint", query.Args{"id": id.String(), "endpoint": endpoint})
}
