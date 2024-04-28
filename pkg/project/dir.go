package project

import "io/fs"

type Directory struct {
	dir   string
	files fs.FS
}

func NewDirectory(dir string, files fs.FS) *Directory {
	return &Directory{dir: dir, files: files}
}

func (p *Directory) Dir() string {
	return p.dir
}

func (p *Directory) Path() string {
	return p.dir
}

func (p *Directory) Files() fs.FS {
	return p.files
}
