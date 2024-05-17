package apps

import (
	"github.com/cryptopunkscc/go-astral-js/mock/appstore"
	"github.com/cryptopunkscc/go-astral-js/target"
)

func List() []target.App {
	return appstore.ListApps()
}
