package test

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/portals"
	"testing"
)

func Test__Builder_find(t *testing.T) {

	scope := feat.Scope[target.Portal]{
		GetPath:      apps.Path,
		TargetFinder: portals.NewFind,
		TargetCache:  target.NewCache[target.Portal](),
	}

	find := scope.GetTargetFind()

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
