//go:build !debug

package main

import (
	"github.com/cryptopunkscc/portal/api/astrald"
	. "github.com/cryptopunkscc/portal/api/target"
	exec2 "github.com/cryptopunkscc/portal/runner/deprecated/exec"
)

func (a *Application) init() {
	a.Resolve = Any[Runnable](
		Skip("node_modules"),
		exec2.Runner{Config: a.Config}.Dist().Try,
		exec2.Runner{Config: a.Config}.DistHost().Try,
		exec2.Runner{Config: a.Config}.Bundle().Try,
		exec2.Runner{Config: a.Config}.BundleHost().Try,
	)
}

func (a *Application) newAstrald() astrald.Runner {
	return &exec2.Astrald{NodeRoot: a.Config.Astrald}
}
