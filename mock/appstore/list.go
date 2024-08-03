package appstore

import (
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/target"
)

func ListApps() []target.App_ {
	return target.List(apps.ResolveAll, portalAppsSource)
}
