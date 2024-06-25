package appstore

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
)

func ListApps() []target.App {
	return apps.FromFS[target.App](portalAppsFs)
}
