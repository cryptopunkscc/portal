package app

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/target"
)

func NewRun(newApi target.NewApi) target.Run[target.App] {
	return Runner{newApi: newApi}.Run
}

type Runner struct{ newApi target.NewApi }

func (r Runner) Run(ctx context.Context, app target.App) (err error) {
	switch v := any(app).(type) {
	case target.AppHtml:
		return wails.NewRunner(r.newApi).Run(ctx, v)
	case target.AppJs:
		return goja.NewRunner(r.newApi).Run(ctx, v)
	default:
		return fmt.Errorf("invalid app target: %v", app.Path())
	}
}
