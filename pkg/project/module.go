package project

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"os"
	"path"
)

type Module struct {
	abs   string
	src   string
	files fs.FS
}

func (m *Module) Abs() string {
	if m.abs != "" {
		return m.abs
	}
	return m.src
}

func (m *Module) Parent() target.Source {
	return NewModule(path.Dir(m.Abs()))
}

func NewModule(src string) (m *Module) {
	m = &Module{}
	m.abs = Abs(src)
	if path.Ext(src) == ".portal" {
		m.src = path.Base(src)
		m.files = os.DirFS(path.Dir(m.abs))
	} else {
		m.src = "."
		m.files = os.DirFS(m.abs)
	}
	return
}

func NewModuleFS(files fs.FS, src string) *Module {
	return &Module{files: files, src: src}
}

func (m *Module) Path() string {
	return m.src
}

func (m *Module) Files() fs.FS {
	return m.files
}

func (m Module) Type() (t target.Type) {
	switch {
	case m.IsFrontend():
		t += target.Frontend
	case m.IsBackend():
		t += target.Backend
	}
	// TODO verify blob type in addition
	if path.Ext(m.src) == ".portal" {
		t += target.Bundle
	}
	return
}

func (m *Module) IsFrontend() bool {
	stat, err := fs.Stat(m.Files(), "index.html")
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func (m *Module) IsBackend() bool {
	stat, err := fs.Stat(m.files, "main.js")
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func (m *Module) Lift() (mm *Module) {
	mm = &(*m)
	if path.Ext(m.Path()) == "" {
		mm.files, _ = fs.Sub(mm.Files(), mm.Path())
		mm.src = "."
	} else {
		mm.files, _ = fs.Sub(mm.Files(), path.Dir(mm.Path()))
		mm.src = path.Base(mm.Path())
	}
	return
}
