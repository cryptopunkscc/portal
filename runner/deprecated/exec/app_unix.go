//go:build !windows

package exec

import (
	"os/exec"
	"syscall"
)

func appKill(cmd *exec.Cmd) {
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
}

func appSysProcArgs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true} // Propagates sig to nested child processes. Required for killing golang apps that are started via go run.
}
