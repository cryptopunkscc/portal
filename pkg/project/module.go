package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
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
		return NewModuleFS(file, os.DirFS(dir))
	}
	return NewModuleFS(src, os.DirFS(src))
}

func NewModuleFS(src string, files fs.FS) *Module {
	return &Module{src: src, files: files}
}

func (m *Module) Path() string {
	return m.src
}

func (m *Module) Files() fs.FS {
	return m.files
}

func (m *Module) Type() runner.Type {
	switch {
	case m.IsFrontend():
		return runner.Frontend
	case m.IsBackend():
		return runner.Backend
	default:
		return runner.Invalid
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
