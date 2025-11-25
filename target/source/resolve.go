package source

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/target"
)

var File target.File = file

func file(path ...string) (t target.Source, err error) {
	abs := target.Abs(path...)
	abs, base := filepath.Split(abs)
	if t = Dir(abs); len(base) > 0 {
		t, err = t.Sub(base)
	}
	return
}

func Dir(path ...string) (dir Dir_) {
	dir = Dir_{}
	dir.absToFS = target.Abs(path...)
	dir.join = filepath.Join
	dir.fs = os.DirFS(dir.absToFS)
	return
}

func Embed(fs fs.FS) (dir Dir_) {
	dir = Dir_{}
	dir.join = path.Join
	dir.fs = fs
	return
}

func FS(fs fs.FS, path ...string) (dir Dir_) {
	dir = Dir_{}
	dir.absToFS = target.Abs(path...)
	dir.join = filepath.Join
	dir.fs = fs
	return
}
