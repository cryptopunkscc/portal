package main

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
	"os"
	"path/filepath"
)

var PortalAppsDir string
var portalAppsSource target.Source

func init() {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	PortalAppsDir = filepath.Join(dir, "portal", "apps")
	err = os.MkdirAll(PortalAppsDir, 0755)
	if err != nil {
		panic(err)
	}
	if portalAppsSource, err = source.File(PortalAppsDir); err != nil {
		panic(err)
	}

}
