package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dist"
	"github.com/cryptopunkscc/go-astral-js/runner/service"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/runner/wails_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/wails_dist"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/msg"
	"reflect"
)

func NewRun(portMsg target.Port) func(newApi target.NewApi) target.Run[target.Portal] {
	return func(newApi target.NewApi) target.Run[target.Portal] {
		return Runner{
			newApi:  newApi,
			portMsg: portMsg,
		}.Run
	}
}

type Runner struct {
	newApi  target.NewApi
	sendMsg target.MsgSend
	portMsg target.Port
}

func (r Runner) Run(ctx context.Context, t target.Portal) (err error) {
	var reloader msg.Reloader
	newApi := func(ctx context.Context, portal target.Portal) target.Api {
		api := r.newApi(ctx, portal)
		handlers := rpc.Handlers{
			r.portMsg.Command: msg.NewHandler(reloader, api).HandleMsg,
		}
		port := r.portMsg.Target(portal).Cmd("")
		service.NewRunner(handlers).Start(ctx, port.String())
		return api
	}
	sendMsg := msg.NewSend(r.portMsg)
	switch v := t.(type) {
	case target.ProjectJs:
		run := goja_dev.NewRunner(newApi, sendMsg)
		reloader = run
		return run.Run(ctx, v)
	case target.ProjectHtml:
		run := wails_dev.NewRunner(newApi) // FIXME propagate sendMsg
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
