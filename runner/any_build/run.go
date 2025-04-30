package any_build

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runner/go_build"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
)

var Run target.Run[string] = dispatcher.Run

var dispatcher = target.Dispatcher{
	Provider: provider,
	Runner:   target.RunSeq,
}

var provider = target.Provider[target.Runnable]{
	Priority: target.Priority{
		target.Match[target.Project_],
		target.Match[target.Dist_],
	},
	Repository: target.Repositories{
		source.Repository,
	},
	Resolve: target.Any[target.Runnable](
		target.Skip("node_modules"),
		npm.BuildRunner().Try,
		go_build.Runner().Try,
		dist.PackRunner.Try,
	),
}
