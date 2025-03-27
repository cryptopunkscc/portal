//go:build !debug

package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
)

func init() {
	application.Resolve = sources.Resolver[Portal_]()
	application.Runners = func(schemaPrefix []string) []Run[Portal_] {
		return []Run[Portal_]{
			app.Runner(exec.DistRun),
			app.Runner(exec.BundleRunner()),
			app.Runner(exec.BundleHostRunner(schemaPrefix...)),
		}
	}
}
