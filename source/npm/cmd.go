package npm

import (
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
)

type BuildNpmAppsOpt struct {
	Clean bool `cli:"clean c"`
	Pack  bool `cli:"pack p"`
}

// BuildNpmApps recursively searches the given path and builds any app bundle it finds.
func BuildNpmApps(opt BuildNpmAppsOpt, path string) (err error) {
	defer plog.TraceErr(&err)
	ref := source.OSRef(path)
	exist := false
	skip := map[string]bool{}
	for _, p := range source.CollectIt(ref, &Project{}) {
		if opt.Clean {
			if err = p.Clean(); err != nil {
				return
			}
		}
		if err = p.Build(); err != nil {
			return
		}
		if opt.Pack {
			if err = p.Project().Pack(); err != nil {
				return
			}
			skip[p.Manifest.Package] = true
		}
		exist = true
	}
	if opt.Pack {
		for _, d := range source.CollectIt(ref, &app.Dist{}) {
			if skip[d.Manifest.Package] {
				continue
			}
			if err = d.Pack(); err != nil {
				return
			}
			exist = true
		}
	}
	if !exist {
		return fs.ErrNotExist
	}
	return
}
