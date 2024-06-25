package dist

import (
	"github.com/cryptopunkscc/portal/target"
)

type source struct {
	target.Source
	manifest *target.Manifest
}

var _ target.Dist = (*source)(nil)

type frontend struct {
	target.Html
	target.Dist
}

type backend struct {
	target.Js
	target.Dist
}

type executable struct {
	target.Exec
	target.Dist
}

var _ target.DistExec = &executable{}

func (m *source) IsApp() {}

func (m *source) IsDist() {}

func (m *source) Type() target.Type {
	return m.Source.Type() + target.TypeDev
}

func (m *source) Manifest() *target.Manifest {
	return m.manifest
}
