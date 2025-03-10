//go:build project

package portald

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
)

var defaultOrder = []int{2, 1, 0}

func (s *Runner[T]) runners(schemaPrefix []string) []Run[Portal_] {
	return []Run[Portal_]{
		app.Runner(exec.DistRun),
		app.Runner(exec.BundleRunner(s.CacheDir)),
		app.Runner(exec.ProjectExecRun),
		app.Runner(exec.ProjectHostRunner(schemaPrefix...)),
	}
}

var resolveApps = sources.Resolver[Portal_]()
