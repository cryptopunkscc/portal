package portal

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

type Dist struct {
	target.Source
	manifest *bundle.Manifest
}

var _ target.Dist = (*Dist)(nil)

func ResolveDist(m target.Source) (module *Dist, err error) {
	sub, err := fs.Sub(m.Files(), m.Path())
	if err != nil {
		return
	}
	manifest, err := bundle.ReadManifestFs(sub)
	if err != nil {
		return
	}
	module = &Dist{Source: m, manifest: &manifest}
	return
}

func (m *Dist) App() {}

func (m *Dist) Dist() {}

func (m *Dist) Type() target.Type {
	return m.Source.Type() + target.TypeDev
}

func (m *Dist) Manifest() *bundle.Manifest {
	return m.manifest
}
