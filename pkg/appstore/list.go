package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	"github.com/cryptopunkscc/go-astral-js/pkg/resolve"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func ListApps() []target.App {
	return array.FromChan(resolve.FromFS[target.App](portalAppsFs))
}
