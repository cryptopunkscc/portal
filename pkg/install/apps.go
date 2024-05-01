package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
)

func (d *PortalDev) BuildJsApps() {
	if len(d.modules) == 0 {
		d.collectPortalLibs()
	}
	if err := project.BuildPortalApps(d.root, "apps", d.modules...); err != nil {
		log.Fatal(err)
	}
}
