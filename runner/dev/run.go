package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/runner/wails_dev"
	"reflect"
)

func NewRunner(bindings target.New) target.Run[target.Portal] {
	return Runner{bindings: bindings, prefix: []string{"dev"}}.Run
}

type Runner struct {
	bindings target.New
	prefix   []string
}

func (r Runner) Run(ctx context.Context, t target.Portal) (err error) {
	prefix := "dev"
	switch v := t.(type) {
	case target.ProjectBackend:
		return goja_dev.NewRunner(r.bindings)(ctx, v)
	case target.ProjectFrontend:
		return wails_dev.NewRunner(r.bindings)(ctx, v)
	case target.AppBackend:
		return goja.NewRunner(r.bindings, prefix)(ctx, v)
	case target.AppFrontend:
		return wails.NewRunner(r.bindings, prefix)(ctx, v)
	default:
		return fmt.Errorf("invalid target %v: %v", reflect.TypeOf(t), t.Path())
	}
}
