package portal

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

type Dist struct {
	target.Source
	manifest *target.Manifest
}

var _ target.Dist = (*Dist)(nil)

type FrontendDist struct {
	target.Frontend
	target.Dist
}

type BackendDist struct {
	target.Backend
	target.Dist
}

var ErrNotDist = errors.New("not a dist")

func ResolveDist(m target.Source) (b target.Dist, err error) {
	if m.IsFile() {
		return nil, ErrNotDist
	}
	m = m.Lift()
	if f, err := m.Files().Open(target.PackageJsonFilename); err == nil {
		_ = f.Close()
		return nil, ErrNotDist
	}
	manifest, err := ReadManifest(m.Files())
	if err != nil {
		return
	}
	b = &Dist{Source: m, manifest: &manifest}
	switch {
	case b.Type().Is(target.TypeFrontend):
		b = &FrontendDist{Dist: b}
	case b.Type().Is(target.TypeBackend):
		b = &BackendDist{Dist: b}
	}
	return
}

func (m *Dist) IsApp() {}

func (m *Dist) IsDist() {}

func (m *Dist) Type() target.Type {
	return m.Source.Type() + target.TypeDev
}

func (m *Dist) Manifest() *target.Manifest {
	return m.manifest
}
