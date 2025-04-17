//go:build !debug

package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runner/exec"
)

func (a *Application[T]) init() {
	a.Resolve = Any[Runnable](
		exec.DistRunner.Try,
		exec.NewBundleRunner(a.Config.Bin).Try,
		exec.NewBundleHostRunner(a.Config.Bin).Try,
	)
}
