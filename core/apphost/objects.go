package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/portal/pkg/resources"
)

func (a *Adapter) Objects() *ObjectsClient {
	return &ObjectsClient{*a.Client, *astrald.NewObjectsClient(a.TargetID, a.Client)}
}

type ObjectsClient struct {
	astrald.Client
	astrald.ObjectsClient
}

func (c *ObjectsClient) Fetch(id *astral.ObjectID, obj astral.Object) (err error) {
	b, err := c.Read(nil, id, 0, 0)
	if err != nil {
		return
	}
	return resources.ReadCanonical(b, obj)
}

func (c *ObjectsClient) Get(id *astral.ObjectID) (obj astral.Object, err error) {
	b, err := c.Read(nil, id, 0, 0)
	if err != nil {
		return
	}
	obj, _, err = astral.DefaultBlueprints.Canonical().Read(b)
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
