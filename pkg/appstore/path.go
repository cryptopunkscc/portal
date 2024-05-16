package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"io/fs"
	"path"
)

func Path(app string) (src string, err error) {
	for t := range project.FindInPath[*project.Bundle](portalAppsDir) {
		m := bundle.Manifest{}
		if err = m.LoadFs(t.Files(), bundle.PortalJson); err != nil {
			return
		}
		if m.Name == app || m.Package == app {
			src = path.Join(portalAppsDir, t.Path())
			return
		}
	}
	err = fs.ErrNotExist
	return
}
