package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"os"
	"os/exec"
)

func Request(executable string) target.Request {
	return func(ctx context.Context, src string, args ...string) (err error) {
		execArgs := append([]string{src}, args...)
		c := exec.CommandContext(ctx, executable, execArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Start()
	}
}
