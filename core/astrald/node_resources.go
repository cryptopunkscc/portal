package astrald

import "github.com/cryptopunkscc/portal/pkg/resources"

func (r *Initializer) initNodeResources() (err error) {
	if r.resources.FileResources == nil {
		r.resources, err = resources.NewFileResources(r.NodeRoot.Get())
	}
	return
}
