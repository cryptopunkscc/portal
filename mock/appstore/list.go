package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
)

func ListApps() []target.App {
	return apps.FromFS[target.App](portalAppsFs)
}
