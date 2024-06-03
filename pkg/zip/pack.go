package zip

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Pack(src string, dst string) (err error) {
	file, err := os.Create(dst)
	if err != nil && errors.Is(err, os.ErrExist) {
		return fmt.Errorf("os.Create: %v", err)
	}
	defer file.Close()
	w := zip.NewWriter(file)
	defer w.Close()

	// copy files to bundle
	if err = filepath.Walk(src, func(p string, d os.FileInfo, err error) error {
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
		stat, err := file.Stat()
		if err != nil {
			return err
		}
		trim, found := strings.CutPrefix(p, src)
		if !found {
			return nil
		}
		h := &zip.FileHeader{
			Name:   trim,
			Method: zip.Deflate,
		}
		h.SetMode(stat.Mode())
		f, err := w.CreateHeader(h)
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
		return nil
	}); err != nil {
		return fmt.Errorf("filepath.Walk: %v", err)
	}
	return
}
