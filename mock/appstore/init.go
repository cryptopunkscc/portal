package appstore

import (
	"io/fs"
	"os"
	"path/filepath"
)

var portalAppsDir string
var portalAppsFs fs.FS

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
	portalAppsFs = os.DirFS(portalAppsDir)
}
