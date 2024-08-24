package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/factory/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
)

func installApps() {
	if err := apps.Default().InstallSources(
		context.TODO(),
		source.Embed(FS),
	); err != nil {
		panic(err)
	}
}
