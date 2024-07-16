package exec

import (
	"os"
	"os/exec"
)

func Run(dir string, cmd ...string) error {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = os.Environ()
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}
