package test

import (
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func ExecCmdRun(t *testing.T, cmd string, args ...string) {
	cc := ExecCmd(cmd, args...)
	err := cc.Run()
	assert.NoError(t, err)
}

func ExecCmd(cmd string, args ...string) *exec.Cmd {
	cc := exec.Command(cmd, args...)
	cc.Stdout = os.Stdout
	cc.Stderr = os.Stderr
	cc.Stdin = os.Stdin
	cc.Env = append(os.Environ())
	cc.Dir = "."
	return cc
}

func BuildInstaller() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		cc := ExecCmd("./mage", "build:installer")
		cc.Dir = "../"
		err := cc.Run()
		assert.NoError(t, err)
	})
}

func PackProject() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		c := ExecCmd("sh", "-c", "git ls-files -co --exclude-standard -z | tar -cf ./test/sources.tar --exclude=./test/sources.tar --null -T -")
		c.Dir = "../"
		err := c.Run()
		assert.NoError(t, err)
	})
}
