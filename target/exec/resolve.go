package exec

import (
	"github.com/cryptopunkscc/portal/api/manifest"
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

func ResolveExec(source target.Source) (exec target.Exec, err error) {
	defer plog.TraceErr(&err)

	p, err := portal.Resolve_(source)
	if err != nil {
		return
	}

	s := Source{}
	var t manifest.Dist
	if err = t.LoadFrom(p.FS()); err != nil {
		return
	}
	s.target = t.Target
	file := s.target.Exec
	stat, err := fs.Stat(p.FS(), file)
	if err != nil {
		return
	}
	if stat.Mode().Perm()&0111 == 0 {
		err = plog.Errorf("not executable %s", file)
		return
	}
	if s.executable, err = p.Sub(file); err != nil {
		return
	}

	exec = s
	return
}

func ResolveProjectExec(source target.Source) (out target.Exec, err error) {
	if out, err = ResolveExec(source); err == nil {
		return
	}

	defer plog.TraceErr(&err)
	p, err := project.Resolve_(source)
	if err != nil {
		return
	}

	if p.Build().Get().Exec == "" {
		err = target.ErrNotTarget
		return
	}

	out = Source{executable: p}
	return
}
