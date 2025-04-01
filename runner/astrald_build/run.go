package astrald_build

import (
	"fmt"
	"github.com/cryptopunkscc/portal/api"
	"os"
	"os/exec"
	"path/filepath"
)

func Run() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot resolve working dir: %v", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot resolve home dir: %v", err)
	}
	astrald := "github.com/cryptopunkscc/astrald@" + api.AstralVersion
	astrald = filepath.Join(home, "go/pkg/mod", astrald)
	out := filepath.Join(wd, "cmd/portal-installer/bin/")
	cmd := exec.Command("go", "build", "-o", out, "./cmd/astrald")
	cmd.Dir = astrald
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
