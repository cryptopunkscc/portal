package bundle

import (
	"archive/zip"
	"encoding/json"
	"github.com/cryptopunkscc/go-astral-js/pkg/build"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Run(dir string) (err error) {
	// build dist if needed
	if _, err = fs.Stat(os.DirFS(dir), "dist"); os.IsNotExist(err) {
		if err = build.Run(dir); err != nil {
			return
		}
	}

	// prepare dist path
	dist := path.Join(dir, "dist")

	// read package.json
	pkg := path.Join(dir, "package.json")
	pkgBytes, err := os.ReadFile(pkg)
	if err != nil {
		return
	}
	p := struct{ Name string }{}
	if err = json.Unmarshal(pkgBytes, &p); err != nil {
		return
	}

	// copy service file
	service := "service.js"
	src := path.Join(dir, "src")
	bytes, err := os.ReadFile(path.Join(src, service))
	if err != nil {
		return
	}
	if err = os.WriteFile(path.Join(dist, service), bytes, 0644); err != nil {
		return
	}

	// create empty bundle
	if err = os.Mkdir(path.Join(dir, "build"), 0775); err != nil && !os.IsExist(err) {
		return
	}
	bundle := path.Join(dir, "build", p.Name+".zip")
	file, err := os.Create(bundle)
	if err != nil {
		return
	}
	defer file.Close()
	w := zip.NewWriter(file)
	defer w.Close()

	// copy files to bundle
	return filepath.Walk(dist, func(p string, d os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		file, err := os.Open(p)
		if err != nil {
			return err
		}

		trim, found := strings.CutPrefix(p, dist)
		if !found {
			return nil
		}
		f, err := w.Create(trim)
		if err != nil {
			return err
		}

		if _, err = io.Copy(f, file); err != nil {
			return err
		}
		if err = w.Flush(); err != nil {
			return err
		}
		_ = file.Close()
		log.Println(trim)
		return nil
	})
}
