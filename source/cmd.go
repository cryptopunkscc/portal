package source

import (
	"io/fs"

	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type BuildNpmAppsOpt struct {
	Clean bool `cli:"clean c"`
	Pack  bool `cli:"pack p"`
}

// BuildNpmApps recursively searches the given path and builds any app bundle it finds.
func BuildNpmApps(opt BuildNpmAppsOpt, path string) (err error) {
	defer plog.TraceErr(&err)
	ref := OSRef(path)
	projects := CollectIt(ref, &NpmProject{})
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

// PublishAppBundles recursively searches the given path and publishes any app bundle it finds.
func PublishAppBundles(path string) (err error) {
	defer plog.TraceErr(&err)
	apps := CollectIt(OSRef(path), &AppBundle{})
	if len(apps) == 0 {
		return fs.ErrNotExist
	}
	objects := astrald.Objects()
	for _, app := range apps {
		if err = app.Publish(objects); err != nil {
			return
		}
	}
	return
}
