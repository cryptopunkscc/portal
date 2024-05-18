package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

type Portal struct {
	target.NodeModule
	manifest *bundle.Manifest
}

var _ target.Project = (*Portal)(nil)

func NewPortalModule(src string) (module *Portal, err error) {
	nodeModule, err := ResolveNodeModule(target.NewModule(src))
	if err != nil {
		return
	}
	return ResolvePortalModule(nodeModule)
}

func ResolvePortalModule(m target.NodeModule) (module *Portal, err error) {
	manifest := bundle.Manifest{}
	sub, err := fs.Sub(m.Files(), m.Path())
	if err != nil {
		return
	}
	if err = manifest.LoadFs(sub, "package.json"); err != nil {
		return
	}
	if err = manifest.LoadFs(sub, bundle.PortalJson); err != nil {
		return
	}
	module = &Portal{NodeModule: m, manifest: &manifest}
	return
}

func (m *Portal) Project() {}

func (m *Portal) Type() target.Type {
	return m.NodeModule.Type() + target.TypeDev
}

func (m *Portal) Manifest() *bundle.Manifest {
	return m.manifest
}
