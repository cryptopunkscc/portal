package bundle

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/objects"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Publisher struct {
	objects.Dir
}

func (p Publisher) Publish(
	bundle target.Bundle_,
) (
	objectID *astral.ObjectID,
	release *Release,
	err error,
) {
	defer plog.TraceErr(&err)

	object := &Object[any]{Bundle_: bundle}
	release = &Release{Release: *bundle.Release()}

	if e, ok := bundle.(target.BundleExec); ok {
		release.Target = e.Runtime().Target()
	}
	if release.ManifestID, err = p.Write(bundle.Manifest()); err != nil {
		return
	}
	if release.BundleID, err = p.Write(object); err != nil {
		return
	}
	if objectID, err = p.Write(release); err != nil {
		return
	}
	return
}
