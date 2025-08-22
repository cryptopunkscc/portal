package util

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

type Cmd struct {
	*exec.Cmd
}

func (c *Cmd) OutputT(t *testing.T) []byte {
	c.Stdout = nil
	output, err := c.Cmd.Output()
	assert.NoError(t, err)
	return output
}

func (c *Cmd) RunT(t *testing.T) {
	err := c.Run()
	assert.NoError(t, err)
}

func (c *Cmd) Defaults() *Cmd {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Env = os.Environ()
	c.Dir = "."
	return c
}

func (c *Cmd) NoStd() *Cmd {
	c.Stdin = nil
	c.Stdout = nil
	c.Stderr = nil
	return c
}

func CommandContext(ctx context.Context, cmd string, args ...string) *Cmd {
	c := &Cmd{exec.CommandContext(ctx, cmd, args...)}
	return c.Defaults()
}

func Command(cmd string, args ...string) *Cmd {
	c := &Cmd{exec.Command(cmd, args...)}
	return c.Defaults()
}
