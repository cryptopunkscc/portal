package project

import (
	"github.com/cryptopunkscc/go-astral-js/target"
)

type source struct {
	target.Source
	manifest *target.Manifest
}

type nodeModule struct {
	source
	target.NodeModule
}

type html struct {
	target.ProjectNodeModule
	target.Html
}

type js struct {
	target.ProjectNodeModule
	target.Js
}

type golang struct {
	source
}

var _ target.Project = (*source)(nil)
var _ target.ProjectNodeModule = (*nodeModule)(nil)
var _ target.ProjectHtml = (*html)(nil)
var _ target.ProjectJs = (*js)(nil)
var _ target.ProjectGo = (*golang)(nil)

func (m *source) IsProject() {}

func (m *source) Manifest() *target.Manifest {
	return m.manifest
}

func (m *source) Dist() (t target.Dist) {
	return Dist[target.Dist](m)
}

func (m *nodeModule) Type() target.Type {
	return m.NodeModule.Type() + target.TypeDev
}

func (m *html) DistHtml() (t target.DistHtml) {
	return Dist[target.DistHtml](m)
}

func (m *js) DistJs() (t target.DistJs) {
	return Dist[target.DistJs](m)
}

func (m *golang) DistGolang() target.DistExec {
	return Dist[target.DistExec](m)
}

func (m *golang) Type() target.Type {
	return target.TypeDev
}
