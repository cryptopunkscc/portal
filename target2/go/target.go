package golang

import (
	. "github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/exec"
	"github.com/cryptopunkscc/portal/target2/portal"
	"io/fs"
)

type project struct {
	Portal[Exec]
}

func (p *project) IsProject() {}
func (p *project) IsGo()      {}
func (p *project) Dist() (t Dist[Exec]) {
	sub, err := p.Sub("dist")
	if err != nil {
		return
	}
	t, err = exec.ResolveDist(sub)
	return
}

func ResolveProject(src Source) (t ProjectGo, err error) {
	p := &project{}
	if _, err = fs.Stat(src.Files(), "main.go"); err != nil {
		return
	}
	if p.Portal, err = portal.Resolve[Exec](src); err != nil {
		return
	}
	t = p
	return
}
