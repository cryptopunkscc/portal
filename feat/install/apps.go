package install

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/runner/dist"
	"github.com/cryptopunkscc/go-astral-js/runner/pack"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
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
