package bundle

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Publisher struct {
	Dir string
}

func (p Publisher) Publish(bundle target.Bundle_) (i *astral.ObjectID, r *Release, err error) {
	defer plog.TraceErr(&err)

	o := &Object[any]{Bundle_: bundle}
	r = &Release{Release: *bundle.Release()}

	if e, ok := bundle.(target.BundleExec); ok {
		r.Target = e.Runtime().Target()
	}

	if r.ManifestID, err = p.publishObject(bundle.Manifest()); err != nil {
		return
	}
	if r.BundleID, err = p.publishObject(o); err != nil {
		return
	}
	if i, err = p.publishObject(r); err != nil {
		return
	}
	return
}

func (p Publisher) publishObject(obj astral.Object) (id *astral.ObjectID, err error) {
	defer plog.TraceErr(&err)

	buf := bytes.NewBuffer(nil)
	if _, err = astral.WriteCanonical(buf, obj); err != nil {
		return
	}
	if id, err = astral.ResolveObjectID(obj); err != nil {
		return
	}

	n := filepath.Join(p.Dir, id.String())
	err = os.WriteFile(n, buf.Bytes(), 0644)
	return
}
