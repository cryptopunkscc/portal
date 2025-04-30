package uninstall

import (
	"fmt"
	"github.com/cryptopunkscc/portal/resolve/app"
	"github.com/cryptopunkscc/portal/resolve/source"
	"log"
	"os"
)

func Runner(dir string) func(id string) (err error) {
	return func(id string) (err error) {
		src := source.Dir(dir)
		for _, t := range app.Resolve_.List(src) {
			if t.Manifest().Match(id) {
				log.Println("Uninstalling", t.Manifest().Package, "from", t.Abs())
				err = os.Remove(t.Abs())
				return
			}
		}
		return fmt.Errorf("%s not found", id)
	}
}
