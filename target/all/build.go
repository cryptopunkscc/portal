package all

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/go"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
)

var BuildRecursive target.Run[string] = buildDispatcher.Run

var buildDispatcher = target.Dispatcher{
	Runner: &target.AsyncRunner{},
	Provider: target.Provider[target.Runnable]{
		Repository: target.Repositories{
			source.Repository,
		},
		Priority: target.Priority{
			target.Match[target.Project_],
			target.Match[target.Dist_],
		},
		Resolve: target.Any[target.Runnable](
			target.Skip("node_modules"),
			npm.BuildRunner().Try,
			golang.BuildRunner().Try,
			dist.PackRunner.Try,
		),
	},
}
