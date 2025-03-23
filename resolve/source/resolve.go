package source

import (
	"embed"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type source struct {
	scheme   string
	external string
	internal string
	files    fs.FS
	join     func(...string) string
	isFile   bool
}

func (s *source) IsDir() bool       { return !s.isFile }
func (s *source) Abs() (abs string) { return s.join(s.external, s.internal) }
func (s *source) Path() string      { return s.internal }
func (s *source) FS() fs.FS         { return s.files }
func (s *source) Sub(src string) (t target.Source, err error) {
	if s.isFile {
		return nil, errors.New("cannot sub file")
	}
	stat, err := fs.Stat(s.files, src)
	if err != nil {
		return
	}
	ts := *s
	if stat.IsDir() {
		if ts.files, err = fs.Sub(s.FS(), src); err != nil {
			return
		}
		ts.external = s.join(s.external, s.internal, src)
		ts.isFile = false
	} else {
		ts.files = s.FS()
		ts.internal = s.join(ts.internal, src)
		ts.isFile = true
	}
	t = &ts
	return
}

func (s *source) File() (fs.File, error) {
	return s.FS().Open(s.Path())
}

func Embed(files embed.FS) target.Source {
	return &source{
		scheme:   "embed",
		files:    files,
		internal: ".",
		join:     path.Join,
	}
}

func FS(files fs.FS, abs string) target.Source {
	return &source{
		scheme:   "files",
		external: abs,
		internal: ".",
		files:    files,
		join:     filepath.Join,
	}
}

func Dir(abs string) target.Source {
	return FS(os.DirFS(abs), abs)
}

func File(path ...string) (t target.Source, err error) {
	abs := target.Abs(filepath.Join(path...))
	abs, file := filepath.Split(abs)
	tt := FS(os.DirFS(abs), abs)
	if file != "" {
		tt, err = tt.Sub(file)
	}
	t = tt
	return
}
