package test

import (
	"context"
	embedApps "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/portal"
	"github.com/cryptopunkscc/portal/target/portals"
	"github.com/cryptopunkscc/portal/target/source"
	"testing"
)

var portalTestCases = []Case[string]{
	{Src: ".", Matchers: []*Target{
		RpcBackend,
		RpcFrontend,
		BasicBackend,
		BasicFrontend,
		ProjectFrontend,
		ProjectBackend,
		ProjectGo,
		Launcher,
		DistExecutable,
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

	cache := target.NewCache[target.Portal]()
	find := portals.Finder.Cached(cache)(findPath, embedFs)

	for _, test := range portalTestCases {
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
