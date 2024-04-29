package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"os"
)

func Uninstall(id string) (err error) {
	for target := range project.BundleTargets(os.DirFS(portalAppsDir), ".") {
		manifest := bundle.Manifest{}
		_ = manifest.LoadFs(target.Files(), bundle.PortalJson)
		if manifest.Package == id {
			err = fs.DeleteFile(target.Path())
			return
		}
	}

	return fmt.Errorf("%s not found", id)
}
