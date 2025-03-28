package main

import (
	"github.com/cryptopunkscc/portal/api/env"
	"os"
	"path/filepath"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	portalDir := filepath.Join(homeDir, ".local/share/portald")
	env.PortaldTokens.SetDir(portalDir, "token")
}
