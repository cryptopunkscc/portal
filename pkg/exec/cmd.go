package exec

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	Ctx    context.Context
	Dir    string
	Cmd    string
	Env    []string
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c Cmd) Default() Cmd {
	c.Env = os.Environ()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func (c Cmd) AddEnv(env ...string) Cmd {
	c.Env = append(c.Env, env...)
	return c
}

func (c Cmd) Parse(cmd string, args ...string) Cmd {
	c.Cmd = cmd
	for _, s := range args {
		chunks := strings.Split(s, " ")
		c.Args = append(c.Args, chunks...)
	}
	return c
}

func (c Cmd) Build() (cmd *exec.Cmd) {
	if c.Ctx != nil {
		cmd = exec.CommandContext(c.Ctx, c.Cmd, c.Args...)
	} else {
		cmd = exec.Command(c.Cmd, c.Args...)
	}
	for _, s := range [][]string{_env, c.Env} {
		cmd.Env = append(cmd.Env, s...)
	}
	cmd.Dir = c.Dir
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return
}
