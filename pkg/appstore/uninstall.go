package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/resolve"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
)

func Uninstall(id string) (err error) {
	for t := range resolve.FromPath[target.Bundle](portalAppsDir) {
		manifest := bundle.Manifest{}
		_ = manifest.LoadFs(t.Files(), bundle.PortalJson)
		if manifest.Name == id || manifest.Package == id {
			log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
			err = fs.DeleteFile(t.Abs())
			return
		}
	}

	return fmt.Errorf("%s not found", id)
}
