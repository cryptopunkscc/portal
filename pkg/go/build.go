package golang

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
	"os/exec"
	"path/filepath"
)

type GoPkgMod struct {
	Url     string
	Version string
}

func (b GoPkgMod) Build(path, out string) (err error) {
	defer plog.TraceErr(&err)
	d, err := b.Path()
	if err != nil {
		return
	}
	c := exec.Command("go", "build", "-o", out, path)
	c.Dir = d
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (b GoPkgMod) Path() (path string, err error) {
	defer plog.TraceErr(&err)
	home, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("cannot resolve home dir: %v", err)
		return
	}
	path = fmt.Sprintf("%s@%s", b.Url, b.Version)
	path = filepath.Join(home, "go/pkg/mod", path)
	return
}
