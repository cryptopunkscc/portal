package install

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/runner/dist"
	"github.com/cryptopunkscc/go-astral-js/runner/pack"
	"log"
)

func (d *PortalDev) buildJsApps() {
	if len(d.modules) == 0 {
		d.collectPortalLibs()
	}
	feat := build.NewFeat(dist.NewRun, pack.Run, d.modules...)
	if err := feat.Dist(context.TODO(), d.root, "apps"); err != nil {
		log.Fatal(err)
	}
}
