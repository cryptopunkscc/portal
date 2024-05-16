package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
)

func Uninstall(id string) (err error) {
	for target := range project.FindInPath[*project.Bundle](portalAppsDir) {
		manifest := bundle.Manifest{}
		_ = manifest.LoadFs(target.Files(), bundle.PortalJson)
		if manifest.Name == id || manifest.Package == id {
			log.Println("Uninstalling", target.Manifest().Package, "from", target.Abs())
			err = fs.DeleteFile(target.Abs())
			return
		}
	}

	return fmt.Errorf("%s not found", id)
}
