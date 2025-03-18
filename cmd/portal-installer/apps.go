package main

import (
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/install"
	"github.com/cryptopunkscc/portal/runtime/dir"
)

func installApps() {
	bundle, err := exec.ResolveBundle(source.Embed(apps.LauncherSvelteFS))
	if err != nil {
		panic(err)
	}
	err = install.Runner{OutputDir: dir.App}.Bundle(bundle)
	if err != nil {
		panic(err)
	}
}
