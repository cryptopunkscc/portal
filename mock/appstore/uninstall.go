package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
	"log"
)

func Uninstall(id string) (err error) {
	for t := range apps.FromPath[target.Bundle](portalAppsDir) {
		manifest, _ := manifest.Read(t.Files())
		if manifest.Name == id || manifest.Package == id {
			log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
			err = fs.DeleteFile(t.Abs())
			return
		}
	}

	return fmt.Errorf("%s not found", id)
}
