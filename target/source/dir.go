package source

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io/fs"
)

type Dir_ struct {
	fs      fs.FS
	join    func(...string) string
	absToFS string
}

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

	p := d.join(src...)

	s, err := fs.Stat(d.fs, p)
	if err != nil {
		return
	}

	if s.IsDir() {
		if d.fs, err = fs.Sub(d.fs, p); err == nil {
			d.absToFS = d.join(d.absToFS, p)
			source = d
		}
		return
	}

	source = File_{
		parent:     d,
		pathToFile: p,
	}
	return
}
