package bundle

import (
	"bytes"
	"encoding/json"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/object"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"os"
	"path/filepath"
)

type Release struct {
	ManifestID *object.ID
	BundleID   *object.ID
	manifest.Release
	manifest.Target
}

func (Release) ObjectType() string { return "app.release" }

func (a *Release) ReadFrom(r io.Reader) (n int64, err error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return
	}
	plog.Println(string(b))
	if err = json.Unmarshal(b, a); err != nil {
		return
	}
	n = int64(len(b))
	return
}

func (a Release) WriteTo(w io.Writer) (n int64, err error) {
	b, err := json.Marshal(a)
	if err != nil {
		return
	}
	nn, err := w.Write(b)
	if err != nil {
		return
	}
	n = int64(nn)
	return
}

type Publisher struct {
	Dir string
}

func (p Publisher) Publish(bundle target.Bundle_) (i *object.ID, r *Release, err error) {
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

func (p Publisher) publishObject(obj astral.Object) (id *object.ID, err error) {
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
