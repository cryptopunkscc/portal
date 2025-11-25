package npm

import (
	"context"
	"os"
	"os/exec"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/deps"
)

func Watch(ctx context.Context, src target.ProjectNpm_) (err error) {
	if err = deps.Check("npm", "-v"); err != nil {
		return
	}
	return npmRunWatch(ctx, src.Abs()).Start()
}

func npmRunWatch(ctx context.Context, src string) *exec.Cmd {
	//cmd := exec.CommandContext(ctx, "gnome-terminal", "--", "npm", "run", "watch")
	cmd := exec.CommandContext(ctx, "npm", "run", "watch")
	cmd.Env = os.Environ()
	cmd.Dir = src
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd
}
