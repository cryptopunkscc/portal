package backend_dev

import (
	"context"
	"os"
	"os/exec"
)

func NpmRunWatch(ctx context.Context, src string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "npm", "run", "watch")
	cmd.Env = os.Environ()
	cmd.Dir = src
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}
