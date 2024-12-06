package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"os"
	"os/exec"
	"strings"
)

func Request(executable string) target.Request {
	return func(ctx context.Context, src string) (err error) {
		args := strings.Split(src, " ")
		c := exec.CommandContext(ctx, executable, args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Start()
	}
}
