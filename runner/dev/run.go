package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dist"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/runner/wails_dev"
	"github.com/cryptopunkscc/go-astral-js/target"
	"reflect"
)

func NewRun(newApi target.NewApi) target.Run[target.Portal] {
	return Runner{newApi: newApi, prefix: []string{"dev"}}.Run
}

type Runner struct {
	newApi target.NewApi
	prefix []string
}

func (r Runner) Run(ctx context.Context, t target.Portal) (err error) {
	prefix := "dev"
	switch v := t.(type) {
	case target.ProjectBackend:
		return goja_dev.NewRun(r.newApi)(ctx, v)
	case target.ProjectFrontend:
		return wails_dev.NewRun(r.newApi)(ctx, v)
	case target.DistBackend:
		return goja_dist.NewRun(r.newApi)(ctx, v)
	case target.AppBackend:
		return goja.NewRun(r.newApi, prefix)(ctx, v)
	case target.AppFrontend:
		return wails.NewRun(r.newApi, prefix)(ctx, v)
	default:
		return fmt.Errorf("invalid target %v: %v", reflect.TypeOf(t), t.Path())
	}
}
