package install

import (
	"github.com/cryptopunkscc/go-astral-js/runner/npm"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
	"log"
	"path"
)

func (d *PortalDev) buildJsLibs() {
	for p := range sources.FromPath[target.NodeModule](path.Join(d.root, "pkg")) {
		if p.PkgJson().IsPortalLib() {
			d.modules = append(d.modules, p)
		}
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

func (d *PortalDev) collectPortalLibs() {
	for p := range sources.FromPath[target.NodeModule](path.Join(d.root, "pkg")) {
		if p.PkgJson().IsPortalLib() {
			d.modules = append(d.modules, p)
		}
	}
}
