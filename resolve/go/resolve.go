package golang

import (
	"encoding/json"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/resolve/exec"
	target2 "github.com/magefile/mage/target"
	"io/fs"
	"path/filepath"
	"slices"
)

type project struct {
	manifest Manifest
	Source
	build Builds
}

func (p *project) Changed(skip ...string) bool  { return GoChanged(p, skip...) }
func (p *project) IsGo()                        {}
func (p *project) Manifest() *Manifest          { return &p.manifest }
func (p *project) MarshalJSON() ([]byte, error) { return json.Marshal(p.Manifest()) }
func (p *project) Build() Builds                { return p.build }
func (p *project) Target() Exec                 { return p.Dist().Target() }
func (p *project) Dist_() Dist_                 { return p.Dist() }
func (p *project) Dist() (t Dist[Exec]) {
	sub, err := p.Sub("dist")
	if err != nil {
		return
	}
	t, err = exec.ResolveDist(sub)
	return
}

var ResolveProject Resolve[ProjectGo] = resolveProject

func resolveProject(src Source) (t ProjectGo, err error) {
	p := &project{}
	if !src.IsDir() {
		return nil, ErrNotTarget
	}
	if _, err = fs.Stat(src.FS(), "main.go"); err != nil {
		return
	}
	if err = all.Unmarshalers.Load(&p.manifest, src.FS(), BuildFilename); err != nil {
		return
	}
	p.Source = src
	p.build = LoadBuilds(src)
	t = p
	return
}

func GoChanged(p Project_, skip ...string) (changed bool) {
	dist_, err := p.Sub("dist")
	if err != nil {
		return true
	}

	dir, err := fs.ReadDir(p.FS(), ".")
	if err != nil {
		panic(err)
	}
	abs := p.Abs()
	skip = append(skip, "build", "dist")
	names := map[string]any{}
	for _, entry := range dir {
		if name := entry.Name(); !slices.Contains(skip, name) {
			names[filepath.Join(abs, name)] = entry
		}
	}

	imports, err := golang.ListImports(abs)
	if err == nil {
		for _, refs := range imports {
			for _, ref := range refs.Refs {
				names[ref] = ref
			}
		}
	}

	var namesArr []string
	for name := range names {
		namesArr = append(namesArr, name)
	}
	if changed, err = target2.Path(dist_.Abs(), namesArr...); err != nil {
		return true
	}
	return
}
