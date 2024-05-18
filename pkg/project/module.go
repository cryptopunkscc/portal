package project

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"os"
	"path"
	"reflect"
)

type Module struct {
	abs   string
	src   string
	files fs.FS
}

func (m *Module) String() string {
	return fmt.Sprintf("%v@%s", reflect.TypeOf(m), m.abs)
}

func (m *Module) Abs() string {
	if m.abs != "" {
		return m.abs
	}
	return m.src
}

func (m *Module) Parent() target.Source {
	dir := path.Dir(m.Abs())
	if path.IsAbs(m.Abs()) {
		return NewModule(dir)
	}
	sub, err := fs.Sub(m.files, dir)
	if err != nil {
		panic(err)
	}
	return NewModuleFS(sub, dir)
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

func NewModuleFS(files fs.FS, src ...string) *Module {
	m := &Module{files: files, src: "."}
	if len(src) > 0 {
		m.src = src[0]
	}
	if len(src) > 1 {
		m.abs = path.Join(src[1:]...)
		if !path.IsAbs(m.abs) {
			println("[WARNING] Module initialized with incorrect absolute path: "+m.abs, m.src)
		}
	}
	return m
}

func (m *Module) Path() string {
	return m.src
}

func (m *Module) Files() fs.FS {
	return m.files
}

func (m *Module) Type() (t target.Type) {
	switch {
	case m.IsFrontend():
		t += target.TypeFrontend
	case m.IsBackend():
		t += target.TypeBackend
	}
	// TODO verify blob type in addition
	if path.Ext(m.src) == ".portal" {
		t += target.TypeBundle
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

func (m *Module) Lift() target.Source {

	// omit if a dir already lifted
	if path.Dir(m.Abs()) == "." {
		return m
	}

	// try lift a dir
	if path.Ext(m.Path()) == "" {
		mm := *m
		mm.files, _ = fs.Sub(m.files, m.src)
		mm.src = "."
		return &mm
	}

	// try lift a file
	if dir := path.Dir(m.src); dir != "" {
		mm := *m
		mm.files, _ = fs.Sub(m.files, path.Dir(m.src))
		mm.src = path.Base(m.src)
		return &mm

	}

	return m
}
