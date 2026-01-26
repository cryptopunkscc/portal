package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	objects "github.com/cryptopunkscc/astrald/mod/objects/client"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/resources"
)

func (a *Adapter) Objects() *ObjectsClient {
	return &ObjectsClient{*a.Client, *objects.New(a.TargetID, a.Client)}
}

type ObjectsClient struct {
	astrald.Client
	ObjectsClient objects.Client
}

func (c *ObjectsClient) Fetch(ctx *astral.Context, id *astral.ObjectID, obj astral.Object) (err error) {
	b, err := c.ObjectsClient.Read(ctx, id, 0, 0)
	if err != nil {
		return
	}
	return resources.ReadCanonical(b, obj)
}

func (c *ObjectsClient) Get(ctx *astral.Context, id *astral.ObjectID) (obj astral.Object, err error) {
	b, err := c.ObjectsClient.Read(ctx, id, 0, 0)
	if err != nil {
		return
	}
	obj, _, err = astral.Decode(b, astral.Canonical())
	return
}

type ScanArgs struct {
	Type   string
	Repo   string
	Out    string
	Follow bool
	Zone   astral.Zone
}

func (c *ObjectsClient) Scan(ctx *astral.Context, repo string, follow bool) (results <-chan *astral.ObjectID, err error) {
	return GoChan[*astral.ObjectID](ctx, c.Client, "objects.scan", query.Args{"repo": repo, "follow": follow})
}

func (c *ObjectsClient) Store(ctx *astral.Context, repo string, object astral.Object) (id *astral.ObjectID, err error) {
	channel, err := c.Client.QueryChannel(ctx, "objects.store", query.Args{"repo": repo})
	if err != nil {
		return
	}
	if err = channel.Send(object); err != nil {
		return
	}
	if object, err = channel.Receive(); err != nil {
		return
	}
	id, ok := object.(*astral.ObjectID)
	if !ok {
		err = plog.Errorf("apphost.ObjectsClient.Store: unexpected object type: %T", object)
	}
	return
}
