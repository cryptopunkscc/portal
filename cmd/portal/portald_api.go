package main

import (
	"context"
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
	"slices"
)

func (a Application) injectPortaldApi(handler *cmd.Handler) {
	if api, err := portald.NewClient().Api(); err == nil {
		a.setupFunctions(api)
		handler.AddSub(api...)
		fixHelp(handler)
	}
}

func (a Application) setupFunctions(handlers cmd.Handlers) {
	for i, handler := range handlers {
		name := handler.Names()[0]
		if handler.Func == "portald" {
			handlers[i].Func = func(ctx context.Context, cmd ...string) (err error) {
				cmd = slices.Insert(cmd, 0, name)
				return a.portaldCli(ctx, cmd...)
			}
		} else {
			handlers[i].Func = func(ctx context.Context, cmd ...string) (err error) {
				cmd = slices.Insert(cmd, 0, name)
				return a.runApp(ctx, nil, cmd)
			}
		}
	}
}

func fixHelp(handler *cmd.Handler) {
	for i, h := range handler.Sub {
		if h.Name == cmd.HelpName {
			handler.Sub[i].Func = cmd.NewHelpFunc(handler)
		}
	}
}
