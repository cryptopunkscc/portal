package project

import "io/fs"

type Directory struct {
	dir   string
	files fs.FS
}

func (p *Directory) Dir() string {
	return p.dir
}

func (p *Directory) Files() fs.FS {
	return p.files
}
