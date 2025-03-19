package dir

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
	"os"
	"path/filepath"
)

func mk(base, name string) (s string) {
	s = filepath.Join(base, name)
	if err := os.MkdirAll(s, 0755); err != nil {
		panic(err)
	}
	return
}

func src(path string) target.Source {
	file, err := source.File(path)
	if err != nil {
		panic(err)
	}
	return file
}
