package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
)

func Uninstall(id string) (err error) {
	targets, err := runner.BundleTargets(portalAppsDir)
	if err != nil {
		return err
	}

	for _, target := range targets {
		manifest := bundle.Manifest{}
		_ = manifest.LoadFs(target.Files, bundle.PortalJson)
		if manifest.Package == id {
			err = fs.DeleteFile(target.Path)
			return
		}
	}

	return fmt.Errorf("%s not found", id)
}
