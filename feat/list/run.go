package list

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func Run() []target.App {
	return appstore.ListApps()
}
