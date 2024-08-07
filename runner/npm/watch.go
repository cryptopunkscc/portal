package npm

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/target"
	"os"
	"os/exec"
)

func RunWatchStart(ctx context.Context, src target.ProjectNpm_) (err error) {
	if err = deps.RequireBinary("npm"); err != nil {
		return
	}
	return cmdRunWatch(ctx, src.Abs()).Start()
}

func cmdRunWatch(ctx context.Context, src string) *exec.Cmd {
	//cmd := exec.CommandContext(ctx, "gnome-terminal", "--", "npm", "run", "watch")
	cmd := exec.CommandContext(ctx, "npm", "run", "watch")
	cmd.Env = os.Environ()
	cmd.Dir = src
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd
}
