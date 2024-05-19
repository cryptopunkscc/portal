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
		return goja_dev.NewBackend(ctx, r.bindings, v).Start()
	case target.ProjectFrontend:
		return wails_dev.NewFrontend(r.bindings, v).Start()
	case target.AppBackend:
		return goja.Run(ctx, r.bindings, v, prefix)
	case target.AppFrontend:
		return wails.Run(r.bindings, v, prefix)
	default:
		return fmt.Errorf("invalid target %v: %v", reflect.TypeOf(t), t.Path())
	}
}
