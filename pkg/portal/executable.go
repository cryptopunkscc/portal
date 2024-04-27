package portal

import (
	"context"
	"os"
	"os/exec"
)

func Executable() string {
	executable, err := os.Executable()
	if err != nil {
		executable = "portal"
	}
	return executable
}

func OpenContext(ctx context.Context, src string) *exec.Cmd {
	c := exec.CommandContext(ctx, Executable(), src)
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
