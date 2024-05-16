package appstore

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

func ListApps() []target.App {
	return array.FromChan(project.FindInFS[target.App](portalAppsFs))
}
