package portal

import (
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

func Open(src string, background bool) (pid int, err error) {
	c := exec.Command(Executable(), src)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
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
