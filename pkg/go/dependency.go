package golang

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Dependency struct {
	Project

	Name    string
	Version string
	Replace string
}

func ProjectDependency(name string) (Dependency, error) { return (&Project{}).Dependency(name) }

func (p *Project) Dependency(name string) (dep Dependency, err error) {
	if len(p.Name) == 0 {
		if err = p.Resolve(); err != nil {
			return
		}
	}
	if dep, err = ParseDependency(string(p.Mod), name); err != nil {
		return
	}
	dep.Project = *p
	return
}

func ParseDependency(goMod string, name string) (dep Dependency, err error) {
	for _, l := range strings.Split(goMod, "\n") {
		if !strings.Contains(l, name) {
			continue
		}

		if !strings.HasPrefix(l, "replace") {
			l = strings.TrimSpace(l)
			c := strings.Split(l, " ")
			name = c[0]
			dep.Name = c[0]
			dep.Version = c[1]
			continue
		}

		c := strings.Split(l, "=>")
		dep.Replace = strings.TrimSpace(c[1])
	}
	if len(dep.Name) == 0 {
		err = plog.Errorf("cannot find dependency: %s", name)
	}
	return
}

func (d *Dependency) Build(pkg, out string) (err error) {
	defer plog.TraceErr(&err)
	p, err := d.Path()
	if err != nil {
		return
	}
	c := exec.Command("go", "build", "-v", "-o", out, pkg)
	c.Dir = p
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (d *Dependency) Path() (path string, err error) {
	if len(d.Replace) > 0 {
		return d.Replace, nil
	}
	defer plog.TraceErr(&err)
	home, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("cannot resolve home dir: %v", err)
		return
	}
	path = fmt.Sprintf("%s@%s", d.Name, d.Version)
	path = filepath.Join(home, "go/pkg/mod", path)
	return
}
