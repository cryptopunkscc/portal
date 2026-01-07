package goja_pro

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/runner/v2/goja/dist"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/js"
)

type Runner struct {
	goja_dist.Runner
	js.Project
}

func NewRunner(core bind.Core) (r *Runner) {
	r = &Runner{}
	r.Core = core
	return
}

func (r Runner) New() source.Source {
	return &r
}

func (r *Runner) Run(ctx context.Context, args ...string) (err error) {
	if err = deps.Check("npm", "-v"); err != nil {
		return
	}
	if err = r.Project.Build(); err != nil {
		return
	}
	if err = r.App.ReadSrc(r.Sub("dist")); err != nil {
		return
	}

	cmd := exec.CommandContext(ctx, "npm", "run", "watch")
	cmd.Env = os.Environ()
	cmd.Dir = r.Path
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err = cmd.Start(); err != nil {
		return
	}

	// Wait 1 sec for npm.Watch finish initial build otherwise runner can restart on the first launch.
	time.Sleep(1 * time.Second)

	r.Runner.App = r.Project.App
	return r.Runner.Run(ctx, args...)
}
