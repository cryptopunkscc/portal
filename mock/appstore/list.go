package appstore

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/apps"
)

func ListApps() []target.App_ {
	return target.List(apps.ResolveAll, portalAppsSource)
}
