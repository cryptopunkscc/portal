//go:build debug

package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runner/exec"
)

func (a *Application[T]) init() {
	a.Order = []int{2, 1, 0}
	a.Resolve = Any[Runnable](
		exec.Runner{Config: a.Config}.Dist().Try,
		exec.Runner{Config: a.Config}.Bundle().Try,
		exec.Runner{Config: a.Config}.Project().Try,
		exec.Runner{Config: a.Config}.ProjectHost().Try,
	)
}
