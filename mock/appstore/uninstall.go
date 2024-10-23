package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"log"
	"os"
)

func Uninstall(id string) (err error) {
	for _, t := range apps.Resolver[target.Bundle_]().List(portalAppsSource) {
		if t.Manifest().Match(id) {
			log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
			err = os.Remove(t.Abs())
			return
		}
	}
	return fmt.Errorf("%s not found", id)
}
