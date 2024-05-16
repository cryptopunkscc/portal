package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

type PortalRawModule struct {
	target.Source
	manifest bundle.Manifest
}

func ResolvePortalRawModule(m target.Source) (module *PortalRawModule, err error) {
	sub, err := fs.Sub(m.Files(), m.Path())
	if err != nil {
		return
	}
	manifest, err := bundle.ReadManifestFs(sub)
	if err != nil {
		return
	}
	module = &PortalRawModule{Source: m, manifest: manifest}
	return
}

func (m *PortalRawModule) App() {}

func (m PortalRawModule) Type() target.Type {
	return m.Source.Type() + target.Dev
}

func (m *PortalRawModule) Manifest() bundle.Manifest {
	return m.manifest
}
