package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/apps"
)

func ListApps() []target.App {
	return array.FromChan(apps.FromFS[target.App](portalAppsFs))
}
