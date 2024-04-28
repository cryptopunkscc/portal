package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
	"os"
)

func (d *PortalDev) appendPortalModules() {
	for p := range project.Find[project.NodeModule](os.DirFS(d.root), "pkg") {
		if p.IsPortalModule() {
			d.modules = append(d.modules, p.Dir())
		}
	}
}

func (d *PortalDev) installJsLibs() {
	for p := range project.Find[project.NodeModule](os.DirFS(d.root), "pkg") {
		if p.IsPortalModule() {
			d.modules = append(d.modules, p.Dir())
		}
		if !p.HasNpmRunBuild() {
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
