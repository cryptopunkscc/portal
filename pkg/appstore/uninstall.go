package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
	"path"
)

func Uninstall(id string) (err error) {
	for target := range project.Find[project.Bundle](portalAppsFs, ".") {
		manifest := bundle.Manifest{}
		_ = manifest.LoadFs(target.Files(), bundle.PortalJson)
		if manifest.Name == id || manifest.Package == id {
			abs := path.Join(portalAppsDir, target.Path())
			log.Println("Uninstalling", target.Manifest().Package, "from", abs)
			err = fs.DeleteFile(abs)
			return
		}
	}

	return fmt.Errorf("%s not found", id)
}
