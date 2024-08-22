package runtime

import (
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/mobile"
	"github.com/cryptopunkscc/portal/api/target"
)

type App struct {
	source  target.App_
	runtime bind.Runtime
}

func (a App) Assets() mobile.Assets { return assets{a.source.Files()} }
func (a App) Runtime() bind.Runtime { return a.runtime }
func (a App) Manifest() *mobile.Manifest {
	manifest := a.source.Manifest()
	return &mobile.Manifest{
		Name:        manifest.Name,
		Title:       manifest.Title,
		Description: manifest.Description,
		Package:     manifest.Package,
		Version:     manifest.Version,
		Icon:        manifest.Icon,
	}
}
