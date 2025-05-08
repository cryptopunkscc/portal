//go:build windows

package wails_pro

import (
	"bytes"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os/exec"
	"strconv"
)

func setParentGID(_ *exec.Cmd) {}

func killProc(cmd *exec.Cmd, devCommand string) {
	// Credit: https://stackoverflow.com/a/44551450
	// For whatever reason, killing an npm script on windows just doesn't exit properly with cancel
	if cmd != nil && cmd.Process != nil {
		kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
		var errorBuffer bytes.Buffer
		var stdoutBuffer bytes.Buffer
		kill.Stderr = &errorBuffer
		kill.Stdout = &stdoutBuffer
		err := kill.Run()
		if err != nil && err.Error() != "exit status 1" {
			println(stdoutBuffer.String())
			println(errorBuffer.String())
			plog.PrintTrace(&err)
		}
	}
}
