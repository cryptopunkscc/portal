package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"path"
)

func Path(app string) (src string, err error) {
	for t := range portal.FromPath[target.Bundle](portalAppsDir) {
		m := target.Manifest{}
		if err = m.LoadFs(t.Files(), target.PortalJson); err != nil {
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
