//go:build !project

package portald

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
)

var defaultOrder = []int{}

func (s *Runner[T]) runners(schemaPrefix []string) (arr []Run[Portal_]) {
	return []Run[Portal_]{
		app.Runner(exec.DistRun),
		app.Runner(exec.BundleRunner(s.CacheDir)),
		app.Runner(exec.BundleHostRunner(s.CacheDir, schemaPrefix...)),
	}
}

var resolveApps = apps.ResolveAll
