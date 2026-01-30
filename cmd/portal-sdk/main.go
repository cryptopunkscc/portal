package main

import (
	"github.com/cryptopunkscc/portal/cmd/portal-sdk/src"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/cryptopunkscc/portal/source/npm"
	"github.com/cryptopunkscc/portal/source/tmpl"
)

func main() { cli.Run(handler) }

var handler = cmd.Handler{
	Name: "portal-dev",
	Desc: "Development kit for Astral apps",
	Sub: cmd.Handlers{
		cmd.Handler{
			Func: tmpl.Create,
			Name: "create c",
			Desc: "Create new Astral app from template",
			Params: cmd.Params{
				{
					Type: "string",
					Desc: "Template name: <svelte|js-rollup|js>",
				},
				{
					Type: "string",
					Desc: "Path to destination directory",
				},
			},
		},
		cmd.Handler{
			Func: npm.BuildNpmApps,
			Name: "build b",
			Desc: "Build Astral apps",
			Params: cmd.Params{
				{
					Name: "c clean",
					Type: "bool",
					Desc: "Clean before build",
				},
				{
					Name: "p pack",
					Type: "bool",
					Desc: "Pack dist directory into app bundle",
				},
				{
					Type: "bool",
					Desc: "Optional path to project directory. Default is '.'",
				},
			},
		},
		cmd.Handler{
			Func: app.Publisher{}.PublishBundles,
			Name: "publish p",
			Desc: "Publish app bundles to Astral.",
			Params: cmd.Params{
				{
					Type: "string",
					Desc: "Optional path to app bundle(s) or containing directory. Default is '.' ",
				},
			},
		},
		cmd.Handler{
			Func: runTarget,
			Name: "run r",
			Desc: "Run HTML/JS app in hot reloading runner.",
			Params: cmd.Params{
				{
					Type: "string",
					Desc: "Path to app bundle, or project containing dev.portal.yml or directory containing portal.json",
				},
			},
		},
		cmd.Handler{
			Name: "user-create",
			Desc: "Create new user and store.",
			Func: portal_sdk.CreateUser,
		},
	},
}
