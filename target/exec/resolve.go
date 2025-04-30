package exec

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/portal"
	"github.com/cryptopunkscc/portal/target/project"
	"io/fs"
)

var ResolveDist = dist.Resolver[target.Exec](ResolveExec)
var ResolveBundle = bundle.Resolver[target.Exec](ResolveDist)
var ResolveProject = project.Resolver[target.Exec](ResolveProjectExec)

func ResolveExec(source target.Source) (t target.Exec, err error) {
	defer plog.TraceErr(&err)

	p, err := portal.Resolve_(source)
	if err != nil {
		return
	}
	defer plog.TraceErr(&err)

	file := p.Manifest().Exec
	stat, err := fs.Stat(p.FS(), file)
	if err != nil {
		return
	}
	if stat.Mode().Perm()&0111 == 0 {
		err = plog.Errorf("not executable %s", file)
		return
	}
	sub, err := p.Sub(file)
	if err != nil {
		return
	}
	t = Source{exec: sub}
	return
}

func ResolveProjectExec(source target.Source) (out target.Exec, err error) {
	defer plog.TraceErr(&err)
	p, err := project.Resolve_(source)
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
	out = Source{exec: p}
	return
}
