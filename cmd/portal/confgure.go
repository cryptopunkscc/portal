package main

import (
	"context"
	"errors"
	"time"

	"github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/config"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (a *Application) Configure() (err error) {
	if len(a.Apphost.Token) > 0 {
		return
	}
	if err = a.resolveConfig(); err != nil {
		return
	}
	if err = a.setupEndpoint(); err != nil {
		return
	}
	if err = a.setupAuthToken(); err != nil {
		return
	}
	return
}

func (a *Application) resolveConfig() (err error) {
	defer plog.TraceErr(&err)
	if err = a.Config.Load(); err != nil {
		if !errors.Is(err, config.ErrNotFound) {
			return // abort when config exist but cannot be loaded for some reason
		}
	}
	if err = a.Config.Build(); err != nil {
		return
	}
	plog.D().Scope("config").Printf("\n%s", a.Config.Yaml())
	return
}

func (a *Application) setupEndpoint() (err error) {
	for _, e := range a.Config.Apphost.Listen {
		a.Apphost.Endpoint = e
		break
	}
	return
}

func (a *Application) setupAuthToken() (err error) {
	defer plog.TraceErr(&err)
	var t *apphost.AccessToken
	tokens := token.Repository{Dir: a.Config.Tokens}
	if t, err = tokens.Get("portal"); err == nil {
		a.Apphost.Token = string(t.Token)
	}
	return
}

func (a *Application) handleConfigurationError(ctx context.Context, err error) error {
	if !errors.Is(err, token.ErrNotCached) {
		return err
	}
	if err = a.startPortald(ctx); err != nil {
		return err
	}
	await := flow.Await{
		UpTo:  5 * time.Second,
		Delay: 50 * time.Millisecond,
		Mod:   6,
		Ctx:   ctx,
	}
	for range await.Chan() {
		if err = a.setupAuthToken(); err == nil {
			break
		}
	}
	return err
}
