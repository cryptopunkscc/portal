package project

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/source"
)

type CreateOpts struct {
	source.CreateOpts
	manifest.Dev
}

func Create(opts CreateOpts) (err error) {
	o := dist.CreateOpts{}
	o.CreateOpts = opts.CreateOpts
	o.Dist = opts.Dist
	if _, err = fs.ReadFile(opts.TemplatesFS, path.Join(opts.Template, "dev")); err == nil {
		o.PortalYml = "dev.portal.yml"
	}
	err = dist.Create(o)
	_ = os.Remove(filepath.Join(opts.Path, "dev"))
	return
}
