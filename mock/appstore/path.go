package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
	"io/fs"
	"path"
)

func Path(app string) (src string, err error) {
	for _, t := range apps.FromPath[target.Bundle](portalAppsDir) {
		var m target.Manifest
		if m, err = manifest.Read(t.Files()); err != nil {
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
