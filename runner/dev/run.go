package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/go_dev"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dev"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails_dev"
	"github.com/cryptopunkscc/portal/runner/wails_dist"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/msg"
	"reflect"
)

func NewRun(
	portMsg target.Port,
	runGo *go_dev.Runner,
) func(newApi target.NewApi) target.Run[target.Portal] {
	return func(newApi target.NewApi) target.Run[target.Portal] {
		return Runner{
			newApi:  newApi,
			portMsg: portMsg,
			runGo:   runGo,
		}.Run
	}
}

type Runner struct {
	newApi  target.NewApi
	portMsg target.Port
	runGo   *go_dev.Runner
}

func (r Runner) Run(ctx context.Context, t target.Portal) (err error) {
	var reloader msg.Reloader
	client := msg.NewClient(r.portMsg)
	newApi := func(ctx context.Context, portal target.Portal) target.Api {
		api := r.newApi(ctx, portal)
		client.Init(reloader, api)
		if err = client.Connect(ctx, t); err != nil {
			plog.Get(ctx).Type(r).P().Println(err)
		}
		return api
	}
	sendMsg := client.Send
	// TODO replace switch with injected factory
	switch v := t.(type) {
	case target.ProjectJs:
		run := goja_dev.NewRunner(newApi, sendMsg)
		reloader = run
		return run.Run(ctx, v)
	case target.ProjectHtml:
		run := wails_dev.NewRunner(newApi) // FIXME propagate sendMsg
		reloader = run
		return run.Run(ctx, v)
	case target.ProjectGo:
		// TODO implement sendMsg support
		plog.Get(ctx).Type(r).Println("running project go")
		run := r.runGo
		reloader = run
		return run.Run(ctx, v)
	case target.DistJs:
		run := goja_dist.NewRunner(newApi, sendMsg)
		reloader = run
		return run.Run(ctx, v)
	case target.DistHtml:
		run := wails_dist.NewRunner(newApi, sendMsg)
		reloader = run
		return run.Run(ctx, v)
	case target.AppJs:
		run := goja.NewRunner(newApi)
		reloader = run
		return run.Run(ctx, v)
	case target.AppHtml:
		run := wails.NewRunner(newApi)
		reloader = run
		return run.Run(ctx, v)
	default:
		return fmt.Errorf("invalid target %v: %v", reflect.TypeOf(t), t.Path())
	}
}
