package exec

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/portal"
	"github.com/cryptopunkscc/portal/target/project"
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
	if !isAnyExecutable(p.FS(), file) {
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

func isAnyExecutable(files fs.FS, file string) bool {
	return IsUnixExecutable(files, file) || IsWindowsExecutable(file)
}

func IsUnixExecutable(files fs.FS, file string) bool {
	if stat, err := fs.Stat(files, file); err != nil {
		return false
	} else {
		return stat.Mode().Perm()&0111 != 0
	}
}

func IsWindowsExecutable(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	pathext := strings.Split(os.Getenv("PATHEXT"), ";")
	for _, e := range pathext {
		if strings.ToLower(e) == ext {
			return true
		}
	}
	return false
}
