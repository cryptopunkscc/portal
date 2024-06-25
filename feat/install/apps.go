package install

import (
	"context"
	"github.com/cryptopunkscc/portal/feat/build"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/pack"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/sources"
	"log"
	"path"
)

func (d *PortalDev) buildJsApps() {
	dir := path.Join(d.root, "target/js/embed/portal")
	libs := sources.FromPath[target.NodeModule](dir)

	feat := build.NewFeat(dist.NewRun, pack.Run, libs...)
	if err := feat.Dist(context.TODO(), d.root, "apps"); err != nil {
		log.Fatal(err)
	}
}
