//go:build !debug

package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runner/exec"
)

func init() {
	application.Resolve = Any[Runnable](
		exec.DistRunner.Try,
		exec.BundleRunner.Try,
		exec.BundleHostRunner.Try,
	)
}
