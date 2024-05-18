package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/runner/wails_dev"
)

type Runner struct {
	bindings target.New
	prefix   []string
}

func NewRunner(bindings target.New) target.Run[target.Portal] {
	return Runner{bindings: bindings, prefix: []string{"dev"}}.Run
}

func (r Runner) Run(ctx context.Context, t target.Portal) (err error) {
	prefix := "dev"
	switch v := t.(type) {
	case target.App:
		typ := v.Type()
		switch {
		case typ.Is(target.TypeBackend):
			return goja.Run(ctx, r.bindings, v, prefix)
		case typ.Is(target.TypeFrontend):
			return wails.Run(r.bindings, v, prefix)
		default:
			return fmt.Errorf("invalid app target: %v", v.Path())
		}
	case target.Project:
		typ := v.Type()
		switch {
		case typ.Is(target.TypeBackend):
			return goja_dev.NewBackend(ctx, r.bindings, v).Start()
		case typ.Is(target.TypeFrontend):
			return wails_dev.NewFrontend(r.bindings, v).Start()
		default:
			return fmt.Errorf("invalid dev target: %v", t.Path())
		}
	}
	return fmt.Errorf("invalid target type %T", t)
}
