package exec

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/base"
	"io/fs"
)

type exec struct{ exec target2.Source }

func (e exec) Executable() target2.Source { return e.exec }

func New(portal target2.Base) (t target2.Exec, err error) {
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

func ResolveExec(portal target2.Source) (t target2.Exec, err error) {
	b, err := base.ResolveBase(portal)
	if err != nil {
		return
	}
	return New(b)
}
