package app

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
)

func NewRunner(bindings target.New) target.Run[target.App] {
	return Runner{bindings: bindings}.Run
}

type Runner struct{ bindings target.New }

func (r Runner) Run(ctx context.Context, app target.App) (err error) {
	switch any(app).(type) {
	case target.AppFrontend:
		return wails.Run(r.bindings, app)
	case target.AppBackend:
		return goja.Run(ctx, r.bindings, app)
	default:
		return fmt.Errorf("invalid app target: %v", app.Path())
	}
}
