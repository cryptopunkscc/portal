package golang

import (
	"github.com/cryptopunkscc/portal/api/target"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	target2 "github.com/magefile/mage/target"
	"io/fs"
	"path/filepath"
	"slices"
)

func Changed(p target.Project_, skip ...string) (changed bool) {
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
