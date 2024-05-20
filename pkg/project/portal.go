package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

type Portal struct {
	target.NodeModule
	manifest *target.Manifest
}

var _ target.Project = (*Portal)(nil)

type Frontend struct {
	target.Project
	target.Frontend
}

type Backend struct {
	target.Project
	target.Backend
}

func NewPortal(src string) (module target.Project, err error) {
	nodeModule, err := ResolveNodeModule(target.NewModule(src))
	if err != nil {
		return
	}
	return ResolvePortal(nodeModule)
}

func ResolvePortal(m target.NodeModule) (b target.Project, err error) {
	manifest := target.Manifest{}
	sub, err := fs.Sub(m.Files(), m.Path())
	if err != nil {
		return
	}
	if err = portal.LoadManifest(&manifest, sub, target.PackageJsonFilename); err != nil {
		return
	}
	if err = portal.LoadManifest(&manifest, sub, target.PortalJsonFilename); err != nil {
		return
	}
	b = &Portal{NodeModule: m, manifest: &manifest}
	switch {
	case b.Type().Is(target.TypeFrontend):
		b = &Frontend{Project: b}
	case b.Type().Is(target.TypeBackend):
		b = &Backend{Project: b}
	}
	return
}

func (m *Portal) IsProject() {}

func (m *Portal) Type() target.Type {
	return m.NodeModule.Type() + target.TypeDev
}

func (m *Portal) Manifest() *target.Manifest {
	return m.manifest
}
