package golang

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"io/fs"
)

type project struct {
	manifest Manifest
	Source
	build Builds
}

func (p *project) IsGo()               {}
func (p *project) Manifest() *Manifest { return &p.manifest }
func (p *project) Build() Builds       { return p.build }
func (p *project) Target() Exec        { return p.Dist().Target() }
func (p *project) Dist_() Dist_        { return p.Dist() }
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
	if !src.IsDir() {
		return nil, ErrNotTarget
	}
	if _, err = fs.Stat(src.Files(), "main.go"); err != nil {
		return
	}
	if err = all.Unmarshalers.Load(&p.manifest, src.Files(), BuildFilename); err != nil {
		return
	}
	p.Source = src
	p.build = LoadBuilds(src)
	t = p
	return
}
