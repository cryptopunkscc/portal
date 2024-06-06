package assets

import (
	_ "embed"
	"io/fs"
	"os"
)

type MapFS struct {
	Map map[string]fs.File
}

func (mfs MapFS) Open(name string) (fs.File, error) {
	if f, b := mfs.Map[name]; b {
		return f, nil
	}
	return nil, os.ErrNotExist
}

type ArrayFs []fs.FS

func (cfs ArrayFs) Open(name string) (fs.File, error) {
	for _, inner := range cfs {
		if file, err := inner.Open(name); err == nil {
			return file, err
		}
	}
	return nil, os.ErrNotExist
}

type OsFS struct{}

func (s OsFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}
