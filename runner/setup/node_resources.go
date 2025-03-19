package setup

import "github.com/cryptopunkscc/portal/pkg/resources"

func (r *Runner) setupResources() (err error) {
	if r.resources.FileResources == nil {
		r.resources, err = resources.NewFileResources(r.NodeRoot)
	}
	return
}
