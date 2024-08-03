package main

import (
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/resolve/source"
)

func installApps() {
	if err := appstore.InstallSource(source.Embed(apps.LauncherSvelteFS)); err != nil {
		panic(err)
	}
}
