package build_all

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/build"
	"github.com/cryptopunkscc/portal/runner/clean"
	"github.com/cryptopunkscc/portal/runner/go_build"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/npm_build"
	"github.com/cryptopunkscc/portal/runner/pack"
	js "github.com/cryptopunkscc/portal/runtime/js/embed"
)

var Run = NewRunner().Run

func NewRunner() *build.Runner {
	return build.NewRunner(
		clean.Runner(),
		multi.Runner[target.Project_](
			go_build.Runner().Portal,
			npm_build.Runner(
				target.Any[target.NodeModule](
					target.Skip("node_modules"),
					target.Try(npm.Resolve)).
					List(source.Embed(js.PortalLibFS))...,
			).Portal,
		),
		pack.Run,
	)
}
