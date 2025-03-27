package astrald

import "github.com/cryptopunkscc/portal/pkg/resources"

func (i *Initializer) initNodeResources() (err error) {
	if i.resources.FileResources == nil {
		i.resources, err = resources.NewFileResources(i.NodeRoot)
	}
	return
}
