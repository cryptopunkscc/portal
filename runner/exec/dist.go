package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/target"
	"os"
	"os/exec"
)

type DistRunner struct {
	ctx    context.Context
	src    target.DistExec
	cancel func() error
}

func NewDistRunner() target.Runner[target.DistExec] {
	return &DistRunner{}
}

func (d *DistRunner) Run(ctx context.Context, src target.DistExec) (err error) {
	d.ctx = ctx
	d.src = src
	return d.Reload()
}

func (d *DistRunner) Reload() error {
	if d.cancel != nil {
		_ = d.cancel()
	}
	cmd := exec.CommandContext(d.ctx, d.src.Executable().Abs())
	d.cancel = cmd.Cancel
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
