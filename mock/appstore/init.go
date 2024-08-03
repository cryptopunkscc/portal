package appstore

import (
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"os"
	"path/filepath"
)

var portalAppsDir string
var portalAppsSource target.Source

func init() {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	portalAppsDir = filepath.Join(dir, "portal", "apps")
	err = os.MkdirAll(portalAppsDir, 0755)
	if err != nil {
		panic(err)
	}
	if portalAppsSource, err = source.File(portalAppsDir); err != nil {
		panic(err)
	}

}
