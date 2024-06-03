package project

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

type source struct {
	target.NodeModule
	manifest *target.Manifest
}

var _ target.ProjectNodeModule = (*source)(nil)

type frontend struct {
	target.ProjectNodeModule
	target.Html
}

type backend struct {
	target.ProjectNodeModule
	target.Js
}

func (m *source) IsProject() {}

func (m *source) Type() target.Type {
	return m.NodeModule.Type() + target.TypeDev
}

func (m *source) Manifest() *target.Manifest {
	return m.manifest
}

func (m *source) Dist() (t target.Dist) {
	return Dist[target.Dist](m)
}

func (m *frontend) DistHtml() (t target.DistHtml) {
	return Dist[target.DistHtml](m)
}

func (m *frontend) DistBackend() (t target.DistJs) {
	return Dist[target.DistJs](m)
}
