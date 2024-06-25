package project

import (
	"github.com/cryptopunkscc/portal/target"
)

var _ target.Project = (*portal)(nil)

type portal struct {
	target.Source
	manifest *target.Manifest
}

func (m *portal) IsProject() {}

func (m *portal) Manifest() *target.Manifest {
	return m.manifest
}

func (m *portal) Dist() (t target.Dist) {
	return Dist[target.Dist](m)
}
