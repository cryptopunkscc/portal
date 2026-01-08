package main

import (
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
		},
		cmd.Handler{
			Func: npm.BuildNpmApps,
			Name: "build b",
			Desc: "Build Astral apps",
		},
		cmd.Handler{
			Func: app.Publisher{}.PublishBundles,
			Name: "publish p",
			Desc: "Publish app bundles to Astral.",
		},
		cmd.Handler{
			Func: runTarget,
			Name: "run r",
			Desc: "Run HTML/JS app in hot reloading runner.",
		},
		cmd.Handler{
			Name: "list l",
			Desc: "List apps and projects recursively found in given path.",
			Func: listTargets,
		},
	},
}
