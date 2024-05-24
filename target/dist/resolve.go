package dist

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
)

var ErrNotDist = errors.New("not a dist target")

func Resolve(src target.Source) (dist target.Dist, err error) {
	if src.IsFile() {
		return nil, ErrNotDist
	}
	src = src.Lift()
	if f, err := src.Files().Open(target.PackageJsonFilename); err == nil {
		_ = f.Close()
		return nil, ErrNotDist
	}
	m, err := manifest.Read(src.Files())
	if err != nil {
		return
	}
	dist = &source{Source: src, manifest: &m}
	switch {
	case dist.Type().Is(target.TypeFrontend):
		dist = &frontend{Dist: dist}
	case dist.Type().Is(target.TypeBackend):
		dist = &backend{Dist: dist}
	}
	return
}
