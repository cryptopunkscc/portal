package app

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/target"
)

func NewRunner(bindings target.NewApi) target.Run[target.App] {
	return Runner{bindings: bindings}.Run
}

type Runner struct{ bindings target.NewApi }

func (r Runner) Run(ctx context.Context, app target.App) (err error) {
	switch v := any(app).(type) {
	case target.AppFrontend:
		return wails.NewRunner(r.bindings)(ctx, v)
	case target.AppBackend:
		return goja.NewRunner(r.bindings)(ctx, v)
	default:
		return fmt.Errorf("invalid app target: %v", app.Path())
	}
}
