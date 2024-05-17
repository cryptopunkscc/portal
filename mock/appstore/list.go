package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
)

func ListApps() []target.App {
	return array.FromChan(apps.FromFS[target.App](portalAppsFs))
}
