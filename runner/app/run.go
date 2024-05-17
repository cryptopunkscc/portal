package app

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
)

type Runner struct {
	bindings target.New
}

func NewRunner(bindings target.New) target.Run[target.App] {
	return Runner{bindings: bindings}.Run
}

func (r Runner) Run(ctx context.Context, app target.App) (err error) {
	typ := app.Type()
	switch {
	case typ.Is(target.Backend):
		return goja.Run(ctx, r.bindings, app)
	case typ.Is(target.Frontend):
		return wails.Run(r.bindings, app)
	default:
		return fmt.Errorf("invalid app target: %v", app.Path())
	}
}
