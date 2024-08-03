package appstore

import (
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
)

func Path(port string) (path string, err error) {
	for _, t := range target.List(
		apps.Resolver[target.Bundle_](),
		portalAppsSource,
	) {
		if t.Manifest().Match(port) {
			path = t.Abs()
			return
		}
	}
	err = fs.ErrNotExist
	return
}
