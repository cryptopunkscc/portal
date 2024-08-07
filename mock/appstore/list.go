package appstore

import (
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/target"
)

func ListApps() []target.App_ {
	return apps.ResolveAll.List(portalAppsSource)
}
