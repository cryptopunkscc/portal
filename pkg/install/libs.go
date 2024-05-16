package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"path"
)

func (d *PortalDev) buildJsLibs() {
	for p := range project.FindInPath[target.NodeModule](path.Join(d.root, "pkg")) {
		if p.IsPortalLib() {
			d.modules = append(d.modules, p)
		}
		if !p.CanNpmRunBuild() {
			continue
		}
		if err := p.NpmInstall(); err != nil {
			log.Fatalln(err)
		}
		if err := p.NpmRunBuild(); err != nil {
			log.Fatalln(err)
		}
	}
}

func (d *PortalDev) collectPortalLibs() {
	for p := range project.FindInPath[target.NodeModule](path.Join(d.root, "pkg")) {
		if p.IsPortalLib() {
			d.modules = append(d.modules, p)
		}
	}
}
