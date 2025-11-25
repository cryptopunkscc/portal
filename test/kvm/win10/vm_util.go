package win10

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func execCmdRun(t *testing.T, cmd string, args ...string) {
	cc := execCmdStd(cmd, args...)
	err := cc.Run()
	assert.NoError(t, err)
}

func execCmdStd(cmd string, args ...string) *exec.Cmd {
	cc := execCmd(cmd, args...)
	cc.Stdout = os.Stdout
	cc.Stderr = os.Stderr
	cc.Stdin = os.Stdin
	return cc
}

func execCmd(cmd string, args ...string) *exec.Cmd {
	cc := exec.Command(cmd, args...)
	cc.Env = append(os.Environ())
	cc.Dir = "."
	return cc
}
