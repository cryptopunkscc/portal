package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"io/fs"
	"os"
	"path"
)

func Path(appPackage string) (src string, err error) {
	for target := range project.Find[project.Bundle](os.DirFS(portalAppsDir), ".") {
		m := bundle.Manifest{}
		if err = m.LoadFs(target.Files(), bundle.PortalJson); err != nil {
			return
		}
		if m.Package == appPackage {
			src = path.Join(portalAppsDir, target.Path())
			return
		}
	}
	err = fs.ErrNotExist
	return
}
