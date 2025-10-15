package bundle

import (
	"encoding/json"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Release struct {
	ManifestID *astral.ObjectID
	BundleID   *astral.ObjectID
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
