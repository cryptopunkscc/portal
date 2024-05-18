package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func ListApps() []target.App {
	return array.FromChan(portal.FromFS[target.App](portalAppsFs))
}
