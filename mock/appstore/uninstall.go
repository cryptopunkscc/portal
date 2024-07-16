package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/manifest"
	"log"
	"os"
)

func Uninstall(id string) (err error) {
	for _, t := range apps.FromPath[target.Bundle](portalAppsDir) {
		m, _ := manifest.Read(t.Files())
		if m.Name == id || m.Package == id {
			log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
			err = os.Remove(t.Abs())
			return
		}
	}

	return fmt.Errorf("%s not found", id)
}
