package exec

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/dist"
	"github.com/cryptopunkscc/portal/resolve/portal"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
)

type exec struct{ exec target.Source }

func (e exec) Executable() target.Source { return e.exec }

func New(portal target.Portal_) (t target.Exec, err error) {
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

func ResolveExec(source target.Source) (t target.Exec, err error) {
	b, err := portal.Resolve_(source)
	if err != nil {
		return
	}
	return New(b)
}

var ResolveDist = dist.Resolver[target.Exec](ResolveExec)
var ResolveBundle = bundle.Resolver[target.Exec](ResolveDist)
