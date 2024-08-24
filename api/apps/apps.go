package apps

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/fsnotify/fsnotify"
	"io"
)

type Apps interface {
	Get(ctx context.Context, pkg string) (target.App_, error)
	List(ctx context.Context) (target.Portals[target.App_], error)
	Observe(ctx context.Context) (apps <-chan App, err error)
	Uninstall(ctx context.Context, pkg string) error
	Install(ctx context.Context, reader io.ReadCloser) error
	InstallSources(ctx context.Context, sources ...target.Source) error
	InstallFromPath(ctx context.Context, path string) error
}

type App struct {
	target.App_ `json:",inline"`
	manifest    *target.Manifest
	Event       *fsnotify.Event `json:"event,omitempty"`
}

func (a *App) Manifest() *target.Manifest {
	if a.App_ != nil {
		return a.App_.Manifest()
	}
	if a.manifest == nil {
		a.manifest = &target.Manifest{}
	}
	return a.manifest
}
