package source

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"io/fs"
	"os"
	"path"
	"strings"
)

func FromPath(src string) target.Source {
	m := &source{}
	m.abs = target.Abs(src)
	if path.Ext(src) == ".portal" {
		m.src = path.Base(src)
		m.files = os.DirFS(path.Dir(m.abs))
	} else {
		m.src = "."
		m.files = os.DirFS(m.abs)
	}
	return m
}

func FromFS(files fs.FS, src ...string) target.Source {
	m := &source{files: files, src: "."}
	if len(src) > 0 {
		m.src = src[0]
	}
	if len(src) > 1 {
		m.abs = path.Join(src[1:]...)
		if !strings.HasSuffix(m.abs, m.src) {
			m.abs = path.Join(m.abs, m.src)
		}
		//if !path.IsAbs(m.abs) {
		//	println("[WARNING] source initialized with incorrect absolute path: "+m.abs, m.src)
		//}
	}
	return m
}
