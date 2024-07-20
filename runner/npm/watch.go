package npm

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"os"
	"os/exec"
)

func CmdRunWatch(ctx context.Context, src string) *exec.Cmd {
	//cmd := exec.CommandContext(ctx, "gnome-terminal", "--", "npm", "run", "watch")
	cmd := exec.CommandContext(ctx, "npm", "run", "watch")
	cmd.Env = os.Environ()
	cmd.Dir = src
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd
}

func RunWatchStart(ctx context.Context, src string) (err error) {
	if err = deps.RequireBinary("npm"); err != nil {
		return
	}
	return CmdRunWatch(ctx, src).Start()
}
