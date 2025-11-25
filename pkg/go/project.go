package golang

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Project struct {
	Dir  string
	Name string
	Mod  []byte
}

func ResolveProject(abs ...string) (project *Project, err error) {
	project = &Project{}
	if err = project.Resolve(abs...); err != nil {
		project = nil
	}
	return
}

func (p *Project) Resolve(abs ...string) (err error) {
	defer plog.TraceErr(&err)
	p.Dir = filepath.Join(abs...)

	if len(p.Dir) == 0 {
		if p.Dir, err = os.Getwd(); err != nil {
			return
		}
	}

	if p.Dir == "." || p.Dir == "" || p.Dir == "/" {
		return plog.Errorf("cannot find root")
	}

	m := filepath.Join(p.Dir, "go.mod")
	if p.Mod, err = os.ReadFile(m); err == nil {
		if len(p.Mod) == 0 {
			return plog.Errorf("empty %s", m)
		}

		p.Name = strings.SplitN(string(p.Mod), "\n", 2)[0]
		p.Name = strings.TrimPrefix(p.Name, "module ")

		if len(p.Name) == 0 {
			return plog.Errorf("cannot parse project name")
		}

		return
	}

	return p.Resolve(filepath.Dir(p.Dir))
}
