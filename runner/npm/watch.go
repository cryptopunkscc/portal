package npm

import (
	"context"
	"os"
	"os/exec"
)

func RunWatch(ctx context.Context, src string) *exec.Cmd {
	//cmd := exec.CommandContext(ctx, "gnome-terminal", "--", "npm", "run", "watch")
	cmd := exec.CommandContext(ctx, "npm", "run", "watch")
	cmd.Env = os.Environ()
	cmd.Dir = src
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd
}
