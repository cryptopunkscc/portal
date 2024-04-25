package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"io/fs"
)

func Path(appPackage string) (src string, err error) {
	targets, err := runner.BundleTargets(portalAppsDir)
	if err != nil {
		return
	}
	for _, target := range targets {
		m := bundle.Manifest{}
		if err = m.LoadFs(target.Files, "portal.json"); err != nil {
			return
		}
		if m.Package == appPackage {
			src = target.Path
			return
		}
	}
	err = fs.ErrNotExist
	return
}
