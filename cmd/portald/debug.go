//go:build debug

package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/portal"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
)

func init() {
	application.Order = []int{2, 1, 0}
	application.Resolve = Any[Portal_](
		Try(exec2.ResolveProject),
		Try(exec2.ResolveDist),
		Try(exec2.ResolveBundle),
		Try(portal.Resolve_),
	)
	application.Runners = func(schemaPrefix []string) []Run[Portal_] {
		return []Run[Portal_]{
			app.Runner(exec.DistRun),
			app.Runner(exec.BundleRunner()),
			app.Runner(exec.ProjectExecRun),
			app.Runner(exec.ProjectHostRunner(schemaPrefix...)),
		}
	}
}
