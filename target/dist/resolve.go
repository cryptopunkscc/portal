package dist

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
)

var ErrNotDist = errors.New("not a dist target")

func Resolve(t target.Source) (d target.Dist, err error) {
	if t.IsFile() {
		return nil, ErrNotDist
	}
	t = t.Lift()
	if f, err := t.Files().Open(target.PackageJsonFilename); err == nil {
		_ = f.Close()
		return nil, ErrNotDist
	}
	m, err := manifest.Read(t.Files())
	if err != nil {
		return
	}
	d = &source{Source: t, manifest: &m}
	switch {
	case d.Type().Is(target.TypeFrontend):
		d = &frontend{Dist: d}
	case d.Type().Is(target.TypeBackend):
		d = &backend{Dist: d}
	}
	return
}
