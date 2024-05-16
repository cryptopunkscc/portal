package install

import (
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"log"
)

func (d *PortalDev) buildJsApps() {
	if len(d.modules) == 0 {
		d.collectPortalLibs()
	}
	if err := build.Dist(d.root, "apps", d.modules...); err != nil {
		log.Fatal(err)
	}
}
