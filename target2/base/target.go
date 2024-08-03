package base

import (
	"github.com/cryptopunkscc/portal/pkg/json"
	"github.com/cryptopunkscc/portal/target"
	. "github.com/cryptopunkscc/portal/target"
)

type base struct {
	Source
	manifest *target.Manifest
}

func (p *base) Manifest() *target.Manifest {
	return p.manifest
}

func (p *base) LoadManifest() error {
	return json.Load(&p.manifest, p.Files(), target.PortalJsonFilename)
}

var ResolveBase Resolve[Base] = resolve

func resolve(src Source) (t Base, err error) {
	p := base{Source: src}
	if err = p.LoadManifest(); err != nil {
		return
	}
	t = &p
	return
}
