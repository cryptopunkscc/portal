package install

import (
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
	"os"
)

func (d *PortalDev) buildJsApps() {
	if len(d.modules) == 0 {
		d.appendPortalModules()
	}
	for p := range project.Find[project.PortalNodeModule](os.DirFS(d.root), "apps") {
		if !p.HasNpmRunBuild() {
			continue
		}
		if err := p.NpmInstall(); err != nil {
			log.Fatal(err)
		}
		if err := p.CopyModules(d.modules); err != nil {
			log.Fatalln(err)
		}
		if err := p.NpmRunBuild(); err != nil {
			log.Fatalln(err)
		}
		if err := build.CopyManifest(p.Dir()); err != nil {
			log.Fatalln(err)
		}
	}
	return
}
