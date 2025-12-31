package main

import (
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/source"
)

func main() { cli.Run(handler) }

var handler = cmd.Handler{
	Name: "portal-dev",
	Desc: "Development kit for Astral apps",
	Sub: cmd.Handlers{
		cmd.Handler{
			Func: source.PublishAppBundles,
			Name: "publish p",
			Desc: "Publish app bundles to Astral",
		},
		cmd.Handler{
			Func: source.BuildNpmApps,
			Name: "build b",
			Desc: "Build Astral apps",
		},
	},
}
