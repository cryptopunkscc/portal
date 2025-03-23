package main

import (
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/install"
)

func installApps() {
	bundle, err := exec.ResolveBundle(source.Embed(apps.Builds))
	if err != nil {
		panic(err)
	}
	r := install.Runner{}
	r.AppsDir = mem.NewVar(env.PortaldApps.MkdirAll())
	r.Tokens.Dir = env.PortaldTokens.MkdirAll()
	err = r.Bundle(bundle)
	if err != nil {
		panic(err)
	}
}
