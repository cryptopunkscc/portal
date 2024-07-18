package exec

import (
	"context"
	"os"
	"os/exec"
)

func Tray(_ context.Context) error {
	c := exec.Command("portal-tray")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}
