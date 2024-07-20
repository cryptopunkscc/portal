package exec

import (
	"os"
	"os/exec"
)

func Run(dir string, cmd ...string) error {
	c := Cmd(dir, cmd...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func Call(dir string, cmd ...string) error {
	c := Cmd(dir, cmd...)
	_, err := c.Output()
	return err
}

func Output(dir string, cmd ...string) ([]byte, error) {
	c := Cmd(dir, cmd...)
	return c.Output()
}

func String(dir string, cmd ...string) (string, error) {
	c := Cmd(dir, cmd...)
	b, err := c.Output()
	return string(b), err
}

func Cmd(dir string, cmd ...string) *exec.Cmd {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = os.Environ()
	c.Dir = dir
	return c
}
