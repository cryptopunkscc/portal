package main

import (
	"context"
	"errors"

	"github.com/cryptopunkscc/portal/pkg/config"
)

type RunArgs struct {
	ConfigPath string `cli:"config c"`
}

func (a *Application) run(ctx context.Context, args RunArgs) (err error) {
	if err = a.start(ctx, args); err != nil {
		return
	}
	return a.Wait()
}

func (a *Application) start(ctx context.Context, args RunArgs) (err error) {
	if err = a.loadConfig(args); err != nil && !errors.Is(err, config.ErrNotFound) {
		return
	}
	if err = a.Configure(); err != nil {
		return
	}
	a.init()
	a.Astrald = a.newAstrald()
	if err = a.Start(ctx); err != nil {
		return
	}
	return
}

func (a *Application) loadConfig(args RunArgs) (err error) {
	var path []string
	if len(args.ConfigPath) > 0 {
		path = append(path, args.ConfigPath)
	}
	err = a.Config.Load(path...)
	return
}
