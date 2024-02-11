package build

import (
	"os"
	"os/exec"
)

func Run(dir string) error {
	cmd := exec.Command("npm", "run", "build")
	cmd.Env = os.Environ()
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
