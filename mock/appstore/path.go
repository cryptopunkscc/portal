package appstore

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"io/fs"
)

func Path(port string) (path string, err error) {
	for _, t := range apps.Resolver[target.Bundle_]().List(portalAppsSource) {
		if t.Manifest().Match(port) {
			path = t.Abs()
			return
		}
	}
	err = fs.ErrNotExist
	return
}
