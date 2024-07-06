package test

import (
	embedApps "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/portal"
	"github.com/cryptopunkscc/portal/target/source"
	"testing"
)

func Test__portal_Resolve_ById__embed_apps(t *testing.T) {

	resolveEmbed := portal.NewResolver[target.App](
		apps.Resolve[target.App](),
		source.FromFS(embedApps.LauncherSvelteFS),
	)

	tests := []Case[string]{
		{Matcher: Launcher, Src: Launcher.Manifest.Name},
		{Matcher: Launcher, Src: Launcher.Manifest.Package},
	}

	for _, test := range tests {
		t.Run(test.Src, func(t *testing.T) {
			app, err := resolveEmbed.Portal(test.Src)
			if err != nil {
				t.Fatal(err)
			}
			test.Assert(t, app)
		})
	}
}
