package goja_pro

import (
	"os"
	"os/exec"
	"time"

	"github.com/cryptopunkscc/portal/pkg/bind/src"
	"github.com/cryptopunkscc/portal/pkg/runner/goja/dist"
	"github.com/cryptopunkscc/portal/pkg/source"
	"github.com/cryptopunkscc/portal/pkg/source/js"
	"github.com/cryptopunkscc/portal/pkg/util/deps"
)

type Runner struct {
	goja_dist.Runner
	js.Project
}

func (r Runner) New() source.Source {
	return &r
}

func (r *Runner) Run(ctx *bind.Core, args ...string) (err error) {
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
