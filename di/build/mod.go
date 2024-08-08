package build

import (
	"github.com/cryptopunkscc/portal/feat/build"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/clean"
	"github.com/cryptopunkscc/portal/runner/go_build"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/npm_build"
	"github.com/cryptopunkscc/portal/runner/pack"
	js "github.com/cryptopunkscc/portal/runtime/js/embed"
	. "github.com/cryptopunkscc/portal/target"
)

func Create() *build.Feat {
	return build.NewFeat(
		clean.NewRunner().Call,
		multi.NewRunner[Project_](
			go_build.NewRun().Portal,
			npm_build.NewRun(
				Any[NodeModule](
					Skip("node_modules"),
					Try(npm.Resolve)).
					List(source.Embed(js.PortalLibFS))...,
			).Portal,
		).Run,
		pack.Run,
	)
}
