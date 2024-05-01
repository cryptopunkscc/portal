package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"os"
	"path"
)

type Module struct {
	src   string
	files fs.FS
}

func NewModule(src string) *Module {
	if path.Ext(src) == ".portal" {
		dir, file := path.Split(src)
		return newModuleFS(file, os.DirFS(dir))
	}
	return newModuleFS(src, os.DirFS(src))
}

func newModuleFS(src string, files fs.FS) *Module {
	return &Module{src: src, files: files}
}

func (m *Module) Path() string {
	return m.src
}

func (m *Module) Files() fs.FS {
	return m.files
}

func (m *Module) Type() target.Type {
	switch {
	case m.IsFrontend():
		return target.Frontend
	case m.IsBackend():
		return target.Backend
	default:
		return target.Invalid
	}
}

func (m *Module) IsFrontend() bool {
	stat, err := fs.Stat(m.files, "index.html")
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func (m *Module) IsBackend() bool {
	_, err := fs.Stat(m.files, "index.html")
	return err != nil
}
