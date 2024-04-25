package appstore

import (
	"os"
	"path"
)

var portalAppsDir string

func init() {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	portalAppsDir = path.Join(dir, "portal", "apps")
	err = os.MkdirAll(portalAppsDir, 0755)
	if err != nil {
		panic(err)
	}
}
