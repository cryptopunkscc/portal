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

type Source struct {
	scheme   string
	external string
	internal string
	files    fs.FS
	join     func(...string) string
	isFile   bool
}

func (s *Source) IsDir() bool       { return !s.isFile }
func (s *Source) Abs() (abs string) { return s.join(s.external, s.internal) }
func (s *Source) Path() string      { return s.internal }
func (s *Source) FS() fs.FS         { return s.files }
func (s *Source) Sub(src string) (t target.Source, err error) {
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

func (s *Source) File() (fs.File, error) {
	return s.FS().Open(s.Path())
}

func Embed(files embed.FS) *Source {
	return &Source{
		scheme:   "embed",
		files:    files,
		internal: ".",
		join:     path.Join,
	}
}

func FS(files fs.FS, abs string) *Source {
	return &Source{
		scheme:   "file",
		external: abs,
		internal: ".",
		files:    files,
		join:     filepath.Join,
	}
}

func Dir(path ...string) *Source {
	abs := target.Abs(path...)
	return FS(os.DirFS(abs), abs)
}

func File(path ...string) (t target.Source, err error) {
	abs, file := filepath.Split(target.Abs(path...))
	if t = Dir(abs); len(file) > 0 {
		t, err = t.Sub(file)
	}
	return
}
