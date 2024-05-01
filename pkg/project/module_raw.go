package project

import "github.com/cryptopunkscc/go-astral-js/pkg/bundle"

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

func (p *PortalRawModule) App() {}

func (p *PortalRawModule) Manifest() bundle.Manifest {
	return p.manifest
}
