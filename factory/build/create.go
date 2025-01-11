package build

import (
	. "github.com/cryptopunkscc/portal/api/target"
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

func Create() *build.Runner {
	return build.NewRunner(
		clean.Runner(),
		multi.Runner[Project_](
			go_build.Runner().Portal,
			npm_build.Runner(
				Any[NodeModule](
					Skip("node_modules"),
					Try(npm.Resolve)).
					List(source.Embed(js.PortalLibFS))...,
			).Portal,
		),
		pack.Run,
	)
}
