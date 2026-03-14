package client

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	objects "github.com/cryptopunkscc/astrald/mod/objects/client"
	"github.com/cryptopunkscc/portal/pkg/util/resources"
)

type Objects struct {
	Astrald astrald.Client
	objects.Client
}

func (c *Objects) Fetch(ctx *astral.Context, id *astral.ObjectID, obj astral.Object) (err error) {
	b, err := c.Read(ctx, id, 0, 0)
	if err != nil {
		return
	}
	return resources.ReadCanonical(b, obj)
}

func (c *Objects) Get(ctx *astral.Context, id *astral.ObjectID) (obj astral.Object, err error) {
	b, err := c.Read(ctx, id, 0, 0)
	if err != nil {
		return
	}
	obj, _, err = astral.Decode(b, astral.Canonical())
	return
}
