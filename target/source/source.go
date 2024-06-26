package source

import (
	"fmt"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
	"log"
	"path"
	"reflect"
)

type source struct {
	abs   string
	src   string
	files fs.FS
}

func (m *source) String() string {
	return fmt.Sprintf("%v@%s", reflect.TypeOf(m), m.abs)
}

func (m *source) Abs() string {
	if m.abs != "" {
		return m.abs
	}
	return m.src
}

func (m *source) Parent() target.Source {
	dir := path.Dir(m.Abs())
	if path.IsAbs(m.Abs()) {
		return FromPath(dir)
	}
	sub, err := fs.Sub(m.files, dir)
	if err != nil {
		panic(err)
	}
	return FromFS(sub, dir)
}

func (m *source) Path() string {
	return m.src
}

func (m *source) Files() fs.FS {
	return m.files
}

func (m *source) IsFile() bool {
	return m.Path() != "." && path.Ext(m.Path()) != ""
}

func (m *source) Type() (t target.Type) {
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

func (m *source) IsFrontend() bool {
	stat, err := fs.Stat(m.Files(), "index.html")
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func (m *source) IsBackend() bool {
	if stat, err := fs.Stat(m.files, "main.js"); err == nil {
		return stat.Mode().IsRegular()
	}
	if stat, err := fs.Stat(m.files, "index.js"); err == nil {
		return stat.Mode().IsRegular()
	}
	return false
}

func (m *source) Lift() target.Source {

	// omit if a dir already lifted
	if m.Path() == "." {
		return m
	}

	stat, err := fs.Stat(m.Files(), m.Path())
	if err != nil {
		log.Printf("[ERROR] cannot lift source %s %s: %v", m.Path(), m.Abs(), err)
		mm := *m
		mm.files = nil
		return &mm
	}
	// try lift a dir
	if stat.IsDir() {
		mm := *m
		var err error
		if mm.files, err = fs.Sub(m.files, m.src); err != nil {
			log.Printf("[ERROR] cannot lift dir %s: %v %T %v", m.src, err, m.files, m.files)
			mm.files = nil
		}
		mm.src = "."
		return &mm
	}

	// try lift a file
	if dir := path.Dir(m.src); dir != "." {
		mm := *m
		var err error
		if mm.files, err = fs.Sub(m.files, path.Dir(m.src)); err != nil {
			log.Printf("[ERROR] cannot lift file %s: %v", m.src, err)
			mm.files = nil
		}
		mm.src = path.Base(m.src)
		return &mm
	}

	return m
}
