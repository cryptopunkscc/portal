package source

import (
	"encoding/json"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type ReleaseMetadata struct {
	ManifestID *astral.ObjectID
	BundleID   *astral.ObjectID
	Release
	Target
}

func init() {
	_ = astral.DefaultBlueprints.Add(
		&ReleaseMetadata{},
	)
}

func (ReleaseMetadata) ObjectType() string { return "app.release" }

func (a *ReleaseMetadata) ReadFrom(r io.Reader) (n int64, err error) {
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

func (a ReleaseMetadata) WriteTo(w io.Writer) (n int64, err error) {
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
