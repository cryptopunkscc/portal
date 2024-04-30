package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"log"
	"os"
)

func (d *PortalDev) buildJsApps() {
	if len(d.modules) == 0 {
		d.collectPortalLibs()
	}
	for m := range project.Find[project.PortalNodeModule](os.DirFS(d.root), "apps") {
		if !m.CanNpmRunBuild() {
			continue
		}
		if err := m.PrepareBuild(d.modules...); err != nil {
			log.Fatal(err)
		}
	}
	return
}
