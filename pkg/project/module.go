package project

import "io/fs"

type Module struct {
	dir   string
	files fs.FS
}

func NewDirectory(dir string, files fs.FS) *Module {
	return &Module{dir: dir, files: files}
}

func (p *Module) Dir() string {
	return p.dir
}

func (p *Module) Path() string {
	return p.dir
}

func (p *Module) Files() fs.FS {
	return p.files
}
