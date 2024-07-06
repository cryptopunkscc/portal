package appstore

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/manifest"
	"io/fs"
	"path"
)

func Path(app string) (src string, err error) {
	for _, t := range apps.FromPath[target.Bundle](portalAppsDir) {
		var m target.Manifest
		if m, err = manifest.Read(t.Files()); err != nil {
			return
		}
		if m.Match(app) {
			src = path.Join(portalAppsDir, t.Path())
			return
		}
	}
	err = fs.ErrNotExist
	return
}
