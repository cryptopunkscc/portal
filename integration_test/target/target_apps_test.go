package test

import (
	"context"
	embedApps "github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/mock/appstore"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/portal"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test__apps_Find__embed_launcher(t *testing.T) {
	resolveEmbed := portal.NewResolver[target.App](
		apps.Resolve[target.App](),
		source.FromFS(embedApps.LauncherSvelteFS),
	)
	findPath := target.Mapper[string, string](
		resolveEmbed.Path,
		appstore.Path,
	)
	find := apps.NewFinder(findPath, embedApps.LauncherSvelteFS).Find

	tests := []Case[string]{
		{Matcher: Launcher, Src: Launcher.Manifest.Name},
		{Matcher: Launcher, Src: Launcher.Manifest.Package},
		{Matcher: Launcher, Src: Launcher.Abs},
		{Matcher: BasicBackend, Src: BasicBackend.Abs},
		{Matcher: BasicFrontend, Src: BasicFrontend.Abs},
		{Matcher: RpcFrontend, Src: RpcFrontend.Abs},
		{Matcher: RpcBackend, Src: RpcBackend.Abs},

		// TODO require building
		//{matcher: ProjectJs, src: ProjectJs.Abs},
		//{matcher: ProjectHtml, src: ProjectHtml.Abs},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Src, func(t *testing.T) {
			apps_, err := find(context.TODO(), test.Src)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, 1, len(apps_))
			for _, app := range apps_ {
				test.Assert(t, app)
			}
		})
	}
}

func Test__apps_Resolve__launcher_from_embed_FS(t *testing.T) {
	dir := Launcher.Abs
	resolve := apps.Resolve[target.App]()
	src := source.FromFS(embedApps.LauncherSvelteFS, dir)

	app, err := resolve(src)
	if err != nil {
		t.Fatal(err)
	}

	Launcher.Assert(t, app)
}
