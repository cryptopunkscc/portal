package main

import (
	"github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (a *Application) Configure() (err error) {
	if len(a.Apphost.AuthToken) > 0 {
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
	if err = a.Config.Load(); err != nil {
		return
	}
	if err = a.Config.Build(); err != nil {
		return
	}
	plog.D().Scope(portal.DefaultConfigFile).Printf("\n%s", a.Config.Yaml())
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
	var t *apphost.AccessToken
	tokens := token.Repository{Dir: a.Config.Tokens}
	if t, err = tokens.Get("portal"); err == nil {
		a.Apphost.AuthToken = string(t.Token)
	}
	return
}
