package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dist"
	"github.com/cryptopunkscc/go-astral-js/runner/service_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/runner/wails_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/wails_dist"
	"github.com/cryptopunkscc/go-astral-js/target"
	"reflect"
)

func NewRun(newApi target.NewApi) target.Run[target.Portal] {
	return Runner{newApi: newApi}.Run
}

type Runner struct {
	newApi target.NewApi
}

func (r Runner) Run(ctx context.Context, t target.Portal) (err error) {
	prefix := "dev"
	ctrlPort := "dev.portal.ctrl"
	var reloader service_dev.Reloader
	newApi := func(ctx context.Context, portal target.Portal) target.Api {
		api := r.newApi(ctx, portal)
		service_dev.NewService(reloader, api).Start(ctx, portal)
		return api
	}
	switch v := t.(type) {
	case target.ProjectBackend:
		run := goja_dev.NewRunner(ctrlPort, newApi)
		reloader = run
		return run.Run(ctx, v)
	case target.ProjectFrontend:
		run := wails_dev.NewRunner(newApi)
		reloader = run
		return run.Run(ctx, v)
	case target.DistBackend:
		run := goja_dist.NewRunner(ctrlPort, newApi)
		reloader = run
		return run.Run(ctx, v)
	case target.DistFrontend:
		run := wails_dist.NewRunner(ctrlPort, newApi)
		reloader = run
		return run.Run(ctx, v)
	case target.AppBackend:
		run := goja.NewRunner(newApi, prefix)
		reloader = run
		return run.Run(ctx, v)
	case target.AppFrontend:
		run := wails.NewRunner(newApi, prefix)
		reloader = run
		return run.Run(ctx, v)
	default:
		return fmt.Errorf("invalid target %v: %v", reflect.TypeOf(t), t.Path())
	}
}
