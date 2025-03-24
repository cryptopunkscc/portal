package main

import (
	"path/filepath"
)

func defaultPortalDir() string { return filepath.Join(localShareDir(), "portald") }
