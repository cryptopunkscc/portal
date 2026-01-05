package npm

import (
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
)

type BuildNpmAppsOpt struct {
	Clean bool `cli:"clean c"`
	Pack  bool `cli:"pack p"`
}

// BuildNpmApps recursively searches the given path and builds any app bundle it finds.
func BuildNpmApps(opt BuildNpmAppsOpt, path string) (err error) {
	defer plog.TraceErr(&err)
	ref := source.OSRef(path)
	projects := source.CollectIt(ref, &Project{})
	if len(projects) == 0 {
		return fs.ErrNotExist
	}
	for _, p := range projects {
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
		}
	}
	return
}
