package main

import (
	"os"
	"os/exec"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func portalRun(cmd ...string) (err error) {
	defer plog.TraceErr(&err)
	c := exec.Command("portal", cmd...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}
