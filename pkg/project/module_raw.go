package project

import "github.com/cryptopunkscc/go-astral-js/pkg/bundle"

type PortalRawModule struct {
	Module
	manifest bundle.Manifest
}

func (p *Module) PortalRawModule() (module *PortalRawModule, err error) {
	manifest, err := bundle.ReadManifestFs(p.files)
	if err != nil {
		return
	}
	module = &PortalRawModule{Module: *p, manifest: manifest}
	return
}
