package apps

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func List() []target.App {
	return appstore.ListApps()
}
