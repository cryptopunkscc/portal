package source

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Dir_ struct {
	fs      fs.FS
	join    func(...string) string
	absToFS string
}

var _ target.Source = Dir_{}

func (d Dir_) Abs() string            { return d.Path() }
func (d Dir_) Path() string           { return d.absToFS }
func (d Dir_) File() (fs.File, error) { return d.fs.Open(".") }
func (d Dir_) IsDir() bool            { return true }
func (d Dir_) FS() fs.FS              { return d.fs }

func (d Dir_) Sub(src ...string) (source target.Source, err error) {
	defer plog.TraceErr(&err)
	if len(src) == 0 {
		source = d
		return
	}

	p := path.Join(src...)

	s, err := fs.Stat(d.fs, p)
	if err != nil {
		return
	}

	if !s.IsDir() {
		source = File_{parent: d, pathToFile: p}
		return
	}

	if d.fs, err = fs.Sub(d.fs, p); err == nil {
		d.absToFS = d.join(d.absToFS, p)
		source = d
	}
	return
}

func (d Dir_) CopyTo(path ...string) (err error) {
	targetDir := filepath.Join(path...)
	return fs.WalkDir(d.FS(), ".", func(path string, entry fs.DirEntry, err error) error {
		defer plog.TraceErr(&err)

		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, filepath.FromSlash(path))

		if entry.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		sub, err := d.Sub(path)
		if err != nil {
			return err
		}

		return sub.CopyTo(targetPath)
	})
}
