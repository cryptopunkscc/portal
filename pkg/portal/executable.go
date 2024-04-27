package portal

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js"
	"os"
	"os/exec"
)

const Name = portal.Name

func Executable() string {
	executable, err := os.Executable()
	if err != nil {
		executable = "portal"
	}
	return executable
}

func OpenContext(ctx context.Context, args ...string) *exec.Cmd {
	c := exec.CommandContext(ctx, Executable(), args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func OpenWithContext(ctx context.Context) func(src string, background bool) (pid int, err error) {
	return func(src string, background bool) (pid int, err error) {
		c := OpenContext(ctx, src)
		if !background {
			err = c.Run()
			return
		}
		if err = c.Start(); err != nil {
			return
		}
		pid = c.Process.Pid
		return
	}
}
