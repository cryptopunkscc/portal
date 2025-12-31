package bundle

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/resources"
)

type Publisher struct {
	Dir
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

type Dir struct {
	Path string
}

func (p Dir) Write(obj astral.Object) (objectID *astral.ObjectID, err error) {
	defer plog.TraceErr(&err)

	buf := bytes.NewBuffer(nil)
	if err = resources.WriteCanonical(buf, obj); err != nil {
		return
	}
	if objectID, err = astral.ResolveObjectID(obj); err != nil {
		return
	}

	n := filepath.Join(p.Path, objectID.String())
	err = os.WriteFile(n, buf.Bytes(), 0644)
	return
}
