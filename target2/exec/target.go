package exec

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/base"
	"io/fs"
)

type exec struct{ exec target.Source }

func (e exec) Executable() target.Source { return e.exec }

func New(portal target.Base) (t target.Exec, err error) {
	file := portal.Manifest().Exec
	stat, err := fs.Stat(portal.Files(), file)
	if err != nil {
		return
	}
	if stat.Mode().Perm()&0111 == 0 {
		err = plog.Errorf("not executable %s", file)
		return
	}
	sub, err := portal.Sub(file)
	if err != nil {
		return
	}
	t = &exec{exec: sub}
	return
}

func ResolveExec(portal target.Source) (t target.Exec, err error) {
	b, err := base.ResolveBase(portal)
	if err != nil {
		return
	}
	return New(b)
}
