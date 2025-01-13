package main

import (
	"context"
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/install"
	runtimeApps "github.com/cryptopunkscc/portal/runtime/apps"
)

func installApps() {
	if err := install.Runner(runtimeApps.Dir).All(
		context.Background(),
		source.Embed(apps.LauncherSvelteFS),
	); err != nil {
		panic(err)
	}
}
