package project

import "github.com/cryptopunkscc/go-astral-js/pkg/bundle"

type PortalNodeModule struct {
	NodeModule
	manifest bundle.Manifest
}

func (m *NodeModule) PortalNodeModule() (module *PortalNodeModule, err error) {
	manifest, err := bundle.ReadManifestFs(m.files)
	if err != nil {
		return
	}
	module = &PortalNodeModule{NodeModule: *m, manifest: manifest}
	return
}
