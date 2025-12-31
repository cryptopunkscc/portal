//go:build debug

package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runner/astrald"
	"github.com/cryptopunkscc/portal/runner/astrald/debug"
	"github.com/cryptopunkscc/portal/runner/deprecated/exec"
)

func (a *Application) init() {
	a.Order = []int{2, 1, 0}
	a.Resolve = Any[Runnable](
		Skip("node_modules"),
		exec.Runner{Config: a.Config}.Dist().Try,
		exec.Runner{Config: a.Config}.Bundle().Try,
		exec.Runner{Config: a.Config}.Project().Try,
		exec.Runner{Config: a.Config}.ProjectHost().Try,
	)
}

func (a *Application) newAstrald() astrald.Runner {
	return &debug.Astrald{
		NodeRoot: a.Config.Astrald,
		DBRoot:   a.Config.AstralDB,
		Ghost:    false,
		Version:  false,
	}
}
