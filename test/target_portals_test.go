package test

import (
	"context"
	embedApps "github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/mock/appstore"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/portal"
	"github.com/cryptopunkscc/go-astral-js/target/portals"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"testing"
)

func Test__portals_Find__embed_launcher(t *testing.T) {
	embedFs := embedApps.LauncherSvelteFS
	resolveEmbed := portal.NewResolver[target.App](
		apps.Resolve[target.App](),
		source.FromFS(embedFs),
	)
	findPath := target.Mapper[string, string](
		resolveEmbed.Path,
		appstore.Path,
	)
	find := target.Cached(portals.NewFind)(findPath, embedFs)

	tests := []Case[string]{
		{Src: ".", Matchers: []*Target{
			RpcBackend,
			RpcFrontend,
			BasicBackend,
			BasicFrontend,
			ProjectFrontend,
			ProjectBackend,
			Launcher,
		}},
		{Src: "test_data/rpc", Matchers: []*Target{
			RpcBackend,
			RpcFrontend,
		}},
		{Matcher: Launcher, Src: Launcher.Manifest.Name},
		{Matcher: Launcher, Src: Launcher.Manifest.Package},
		{Matcher: Launcher, Src: Launcher.Abs},
		{Matcher: BasicBackend, Src: BasicBackend.Abs},
		{Matcher: BasicFrontend, Src: BasicFrontend.Abs},
		{Matcher: RpcFrontend, Src: RpcFrontend.Abs},
		{Matcher: RpcBackend, Src: RpcBackend.Abs},
		{Matcher: ProjectBackend, Src: ProjectBackend.Abs},
		{Matcher: ProjectFrontend, Src: ProjectFrontend.Abs},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Src, func(t *testing.T) {
			apps_, err := find(context.TODO(), test.Src)
			if err != nil {
				t.Fatal(err)
			}
			for _, app := range apps_ {
				test.Assert(t, app)
			}
		})
	}
}
