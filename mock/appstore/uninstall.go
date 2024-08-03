package appstore

import (
	"fmt"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"os"
)

func Uninstall(id string) (err error) {
	for _, t := range target.List(
		apps.Resolver[target.Bundle_](),
		portalAppsSource,
	) {
		if t.Manifest().Match(id) {
			log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
			err = os.Remove(t.Abs())
			return
		}
	}
	return fmt.Errorf("%s not found", id)
}
