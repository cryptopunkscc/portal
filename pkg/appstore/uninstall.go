package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
)

func Uninstall(id string) (err error) {
	for t := range portal.FromPath[target.Bundle](portalAppsDir) {
		manifest := target.Manifest{}
		_ = manifest.LoadFs(t.Files(), target.PortalJsonFilename)
		if manifest.Name == id || manifest.Package == id {
			log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
			err = fs.DeleteFile(t.Abs())
			return
		}
	}

	return fmt.Errorf("%s not found", id)
}
