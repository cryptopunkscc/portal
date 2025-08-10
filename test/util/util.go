package util

import (
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

func (c *Cmd) NoStd() *Cmd {
	c.Stdin = nil
	c.Stdout = nil
	c.Stderr = nil
	return c
}

func Command(cmd string, args ...string) *Cmd {
	cc := exec.Command(cmd, args...)
	cc.Stdout = os.Stdout
	cc.Stderr = os.Stderr
	cc.Stdin = os.Stdin
	cc.Env = append(os.Environ())
	cc.Dir = "."
	return &Cmd{cc}
}
