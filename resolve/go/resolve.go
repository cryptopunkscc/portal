package golang

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"io/fs"
)

var ResolveProject target.Resolve[target.ProjectGo] = resolveProject

func resolveProject(src target.Source) (t target.ProjectGo, err error) {
	defer plog.TraceErr(&err)

	t, ok := src.(target.ProjectGo)
	if ok {
		return
	}
	if !src.IsDir() {
		return nil, target.ErrNotTarget
	}
	if _, err = fs.Stat(src.FS(), "main.go"); err != nil {
		return
	}

	p, err := exec.ResolveProject(src)
	if err != nil {
		return
	}

	t = &Source{p}
	return
}
