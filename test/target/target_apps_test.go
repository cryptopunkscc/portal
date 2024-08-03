package test

import (
	embedApps "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"testing"
)

func Test__apps_Resolve__launcher_from_embed_FS(t *testing.T) {
	src, err := source.Embed(embedApps.LauncherSvelteFS).Sub(Launcher.Abs)
	if err != nil {
		t.Fatal(err)
	}

	app, err := apps.ResolveAll(src)
	if err != nil {
		t.Fatal(err)
	}

	Launcher.Assert(t, app)
}
