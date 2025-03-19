package exec

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/dist"
	"github.com/cryptopunkscc/portal/resolve/portal"
	"github.com/cryptopunkscc/portal/resolve/project"
	"github.com/cryptopunkscc/portal/resolve/unknown"
	"io/fs"
)

type exec struct{ exec target.Source }

func (e exec) Executable() target.Source { return e.exec }

func New(portal target.Portal_) (t target.Exec, err error) {
	file := portal.Manifest().Exec
	stat, err := fs.Stat(portal.FS(), file)
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

func ResolveProjectExec(source target.Source) (out target.Exec, err error) {
	p, err := unknown.ResolveProject(source)
	if err != nil {
		return
	}
	m := p.Manifest()
	if m.Exec == "" {
		e := target.GetBuild(p).Exec
		if e == "" {
			err = plog.Errorf("exec not specified for %s", source.Abs())
			return
		}
		m.Exec = e
	}
	out = exec{exec: p}
	return
}

var _ target.Exec = exec{}

var ResolveDist = dist.Resolver[target.Exec](ResolveExec)
var ResolveBundle = bundle.Resolver[target.Exec](ResolveDist)
var ResolveProject = project.Resolver[target.Exec](ResolveProjectExec)
