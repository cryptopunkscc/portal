package desktop

import (
	"github.com/cryptopunkscc/portal/api/env"
	"path/filepath"
)

func init() {
	env.PortaldHome.Default(defaultPortalHome)
	env.AstraldHome.Default(defaultAstraldHome)
	env.AstraldDb.Default(defaultAstraldHome)
	env.PortaldTokens.Default(defaultTokensDir)
	env.PortaldApps.Default(defaultAppsDir)
}

func defaultPortalHome() string  { return filepath.Join(userConfigDir(), "portald") }
func defaultAstraldHome() string { return filepath.Join(env.PortaldHome.Get(), "astrald") }
func defaultTokensDir() string   { return filepath.Join(env.PortaldHome.Get(), "tokens") }
func defaultAppsDir() string     { return filepath.Join(env.PortaldHome.Get(), "apps") }
