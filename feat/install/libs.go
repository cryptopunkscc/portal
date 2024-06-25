package install

import (
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/sources"
	"log"
	"path"
)

func (d *PortalDev) buildJsLibs() {
	dir := path.Join(d.root, "target/js")
	libs := sources.FromPath[target.NodeModule](dir)
	for _, p := range libs {
		if !p.PkgJson().CanBuild() {
			continue
		}
		if err := npm.Install(p); err != nil {
			log.Fatalln(err)
		}
		if err := npm.RunBuild(p); err != nil {
			log.Fatalln(err)
		}
	}
}
