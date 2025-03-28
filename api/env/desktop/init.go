package desktop

import (
	"github.com/cryptopunkscc/portal/api/env"
	"path/filepath"
)

func init() {
	env.PortaldHome.Default(defaultPortalHome)
	env.AstraldHome.Default(defaultAstraldHome)
	env.AstraldDb.Default(defaultAstraldHome)
}

func defaultPortalHome() string  { return filepath.Join(userConfigDir(), "portald") }
func defaultAstraldHome() string { return filepath.Join(userConfigDir(), "astrald") }
