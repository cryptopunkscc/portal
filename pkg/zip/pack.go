package zip

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func Pack(src, dst string) (err error) {
	defer plog.TraceErr(&err)
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		err = PackFS(os.DirFS(src), ".", dst)
	} else {
		dir, file := filepath.Split(src)
		err = PackFS(os.DirFS(dir), file, dst)
	}
	return
}

func PackFS(files fs.FS, src, dst string) (err error) {
	defer plog.TraceErr(&err)
	f, err := os.Create(dst)
	if err != nil && errors.Is(err, os.ErrExist) {
		return
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	// copy files to bundle
	return fs.WalkDir(files, src, func(p string, e fs.DirEntry, err error) error {
		defer plog.TraceErr(&err)
		if err != nil {
			return err
		}
		if e.IsDir() {
			return nil
		}
		f, err := files.Open(p)
		if err != nil {
			return err
		}
		s, err := f.Stat()
		if err != nil {
			return err
		}
		p = strings.TrimPrefix(p, src)
		h := &zip.FileHeader{
			Name:   p,
			Method: zip.Deflate,
		}
		h.SetMode(s.Mode())
		hw, err := w.CreateHeader(h)
		if err != nil {
			return err
		}

		if _, err = io.Copy(hw, f); err != nil {
			return err
		}
		if err = w.Flush(); err != nil {
			return err
		}
		_ = f.Close()
		return nil
	})
}
