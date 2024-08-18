package appstore

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
)

func ListApps() []target.App_ {
	return apps.ResolveAll.List(portalAppsSource)
}
