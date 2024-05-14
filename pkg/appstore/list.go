package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/arr"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func ListApps() []target.App {
	return arr.FromChan(project.Find[target.App](portalAppsFs, "."))
}
