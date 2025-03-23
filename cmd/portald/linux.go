package main

import (
	"github.com/cryptopunkscc/portal/core/env"
	"os"
	"path/filepath"
)

func init() {
	// portald
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	portalDir := filepath.Join(homeDir, ".local/share/portald")
	env.PortaldTokens.SetDir(portalDir, "token")
	env.PortaldBin.SetDir(portalDir, "token")

	// astrald
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	env.AstraldHome.SetDir(configDir, "astrald")
	env.AstraldDb.SetDir(configDir, "astrald")
}
