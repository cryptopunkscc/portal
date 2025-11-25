package target

import (
	"io/fs"
	"path/filepath"
	"slices"

	target2 "github.com/magefile/mage/target"
)

func Changed(p Project_, skip ...string) (changed bool) {
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
	var names []string
	for _, entry := range dir {
		if slices.Contains(skip, entry.Name()) {
			continue
		}
		names = append(names, filepath.Join(abs, entry.Name()))
	}
	if changed, err = target2.Path(dist_.Abs(), names...); err != nil {
		return true
	}
	return
}
