package project

import "github.com/cryptopunkscc/go-astral-js/pkg/bundle"

type PortalRawModule struct {
	Directory
	manifest bundle.Manifest
}

func (p *Directory) PortalRawModule() (module *PortalRawModule, err error) {
	manifest, err := bundle.ReadManifestFs(p.files)
	if err != nil {
		return
	}
	module = &PortalRawModule{Directory: *p, manifest: manifest}
	return
}
