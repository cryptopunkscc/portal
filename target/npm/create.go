package npm

import (
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/target/project"
	"github.com/cryptopunkscc/portal/target/source"
)

type CreateOpts struct {
	source.CreateOpts
	manifest.Dist
}

func Create(opts CreateOpts) (err error) {
	o := project.CreateOpts{}
	o.CreateOpts = opts.CreateOpts
	o.Path = opts.Path
	o.Dist = opts.Dist
	return project.Create(o)
}
