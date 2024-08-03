package apps

import (
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/target"
)

func List() []target.App_ {
	return appstore.ListApps()
}
