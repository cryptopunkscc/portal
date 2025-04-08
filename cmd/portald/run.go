package main

import (
	"context"
	"github.com/cryptopunkscc/portal/runner/exec"
)

type RunArgs struct {
	ConfigPath string `cli:"config c"`
}

func (a *Application[T]) run(ctx context.Context, args RunArgs) (err error) {
	if err = a.start(ctx, args); err != nil {
		return
	}
	return a.Wait()
}

func (a *Application[T]) start(ctx context.Context, args RunArgs) (err error) {
	if err = a.loadConfig(args); err != nil {
		return
	}
	if err = a.Configure(); err != nil {
		return
	}
	a.init()
	a.Astrald = &exec.Astrald{NodeRoot: a.Config.Astrald}
	if err = a.Start(ctx); err != nil {
		return
	}
	return
}

func (a *Application[T]) loadConfig(args RunArgs) (err error) {
	return a.Config.Load(args.ConfigPath)
}
