package project

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

type source struct {
	target.NodeModule
	manifest *target.Manifest
}

var _ target.Project = (*source)(nil)

type frontend struct {
	target.Project
	target.Frontend
}

type backend struct {
	target.Project
	target.Backend
}

func (m *source) IsProject() {}

func (m *source) Type() target.Type {
	return m.NodeModule.Type() + target.TypeDev
}

func (m *source) Manifest() *target.Manifest {
	return m.manifest
}
