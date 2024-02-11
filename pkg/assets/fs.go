package assets

import (
	"archive/zip"
	_ "embed"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"io/fs"
	"os"
)

func BundleFS(bundleType string, path string) (f fs.FS, err error) {
	switch bundleType {
	case TypeDir:
		f = os.DirFS(path)
	case TypeHtml:
		f, err = SingleFileFs(path, "index.html")
	case TypeJs:
		f, err = SingleFileFs(path, "service.js")
	case TypeZip:
		f, err = zip.OpenReader(path)
	}
	if err == nil {
		f = ArrayFs{[]fs.FS{f, apphost.JsFs()}}
	}
	return
}

func SingleFileFs(path string, name string) (f fs.FS, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	return MapFS{map[string]fs.File{name: file}}, err
}

type MapFS struct {
	Map map[string]fs.File
}

func (mfs MapFS) Open(name string) (fs.File, error) {
	if f, b := mfs.Map[name]; b {
		return f, nil
	}
	return nil, os.ErrNotExist
}

type ArrayFs struct {
	Array []fs.FS
}

func (cfs ArrayFs) Open(name string) (fs.File, error) {
	for _, inner := range cfs.Array {
		if file, err := inner.Open(name); err == nil {
			return file, err
		}
	}
	return nil, os.ErrNotExist
}
