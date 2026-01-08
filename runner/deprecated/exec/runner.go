package exec

import (
	"context"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/portal"
)

type Runner struct {
	portal.Config
	RunAppFunc RunApp
}

func DefaultRunner() (r Runner) {
	if err := r.Build(); err != nil {
		panic(err)
	}
	return
}

func (r Runner) RunApp(ctx context.Context, manifest manifest.App, path string, args ...string) (err error) {
	if r.RunAppFunc != nil {
		return r.RunAppFunc(ctx, manifest, path, args...)
	}
	return AppRunner{r.Config}.RunApp(ctx, manifest, path, args...)
}
