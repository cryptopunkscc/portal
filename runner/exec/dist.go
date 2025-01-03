package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"os"
	"os/exec"
)

type dist struct {
	ctx    context.Context
	src    target.DistExec
	cancel func() error
}

func Dist() target.Runner[target.DistExec] {
	return &dist{}
}

func (d *dist) Run(ctx context.Context, src target.DistExec) (err error) {
	d.ctx = ctx
	d.src = src
	return d.Reload()
}

func (d *dist) Reload() error {
	if d.cancel != nil {
		_ = d.cancel()
	}
	abs := d.src.Target().Executable().Abs()
	cmd := exec.CommandContext(d.ctx, abs)
	d.cancel = cmd.Cancel
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
