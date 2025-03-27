package main

import (
	"path/filepath"
)

func defaultPortalHome() string { return filepath.Join(localShareDir(), "portald") }
