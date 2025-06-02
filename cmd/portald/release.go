//go:build !debug

package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runner/exec"
)

func (a *Application[T]) init() {
	a.Resolve = Any[Runnable](
		Skip("node_modules"),
		exec.Runner{Config: a.Config}.Dist().Try,
		exec.Runner{Config: a.Config}.Bundle().Try,
		exec.Runner{Config: a.Config}.BundleHost().Try,
	)
}
