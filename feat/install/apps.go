package install

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"log"
)

func (d *PortalDev) buildJsApps() {
	if len(d.modules) == 0 {
		d.collectPortalLibs()
	}
	feat := build.NewFeat(d.modules...)
	if err := feat.Dist(context.TODO(), d.root, "apps"); err != nil {
		log.Fatal(err)
	}
}
