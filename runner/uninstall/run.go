package uninstall

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"log"
	"os"
)

func Runner(dir mem.String) func(id string) (err error) {
	return func(id string) (err error) {
		src := source.Dir(dir.Require())
		for _, t := range apps.Resolver[target.Bundle_]().List(src) {
			if t.Manifest().Match(id) {
				log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
				err = os.Remove(t.Abs())
				return
			}
		}
		return fmt.Errorf("%s not found", id)
	}
}
