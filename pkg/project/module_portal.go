package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

var _ target.Project = &PortalNodeModule{}

type PortalNodeModule struct {
	target.NodeModule
	manifest *bundle.Manifest
}

func NewPortalNodeModule(src string) (module *PortalNodeModule, err error) {
	nodeModule, err := ResolveNodeModule(NewModule(src))
	if err != nil {
		return
	}
	return ResolvePortalNodeModule(nodeModule)
}

func ResolvePortalNodeModule(m target.NodeModule) (module *PortalNodeModule, err error) {
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
	module = &PortalNodeModule{NodeModule: m, manifest: &manifest}
	return
}

func (m PortalNodeModule) Type() target.Type {
	return m.NodeModule.Type() + target.Dev
}

func (m *PortalNodeModule) Manifest() *bundle.Manifest {
	return m.manifest
}
