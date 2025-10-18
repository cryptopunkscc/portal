package exec

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

func appKill(cmd *exec.Cmd) {
	_ = windows.GenerateConsoleCtrlEvent(windows.CTRL_C_EVENT, uint32(cmd.Process.Pid))
}

func appSysProcArgs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}
