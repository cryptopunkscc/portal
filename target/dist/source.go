package dist

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

type source struct {
	target.Source
	manifest *target.Manifest
}

var _ target.Dist = (*source)(nil)

type frontend struct {
	target.Frontend
	target.Dist
}

type backend struct {
	target.Backend
	target.Dist
}

func (m *source) IsApp() {}

func (m *source) IsDist() {}

func (m *source) Type() target.Type {
	return m.Source.Type() + target.TypeDev
}

func (m *source) Manifest() *target.Manifest {
	return m.manifest
}
