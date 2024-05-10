package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

type PortalRawModule struct {
	*Module
	manifest bundle.Manifest
}

func (m *Module) PortalRawModule() (module *PortalRawModule, err error) {
	manifest, err := bundle.ReadManifestFs(m.files)
	if err != nil {
		return
	}
	module = &PortalRawModule{Module: m, manifest: manifest}
	return
}

func (m *PortalRawModule) App() {}

func (m PortalRawModule) Type() target.Type {
	return m.Module.Type() + target.Dev
}

func (m *PortalRawModule) Manifest() bundle.Manifest {
	return m.manifest
}
