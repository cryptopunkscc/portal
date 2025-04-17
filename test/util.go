package test

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func CleanDir(t *testing.T, path ...string) string {
	Clean(path...)
	return Dir(t, path...)
}

func Dir(t *testing.T, path ...string) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := ".test"
	if len(path) > 0 {
		dir = filepath.Join(path...)
	}
	return filepath.Join(wd, dir)
}

func Mkdir(t *testing.T, path ...string) (d string) {
	d = Dir(t, path...)
	if err := os.MkdirAll(d, 0755); err != nil {
		t.Fatal(err)
	}
	return
}

func Assert[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func Copy(src target.Source, path ...string) target.Source {
	out := "test_data"
	if len(path) > 0 {
		out = filepath.Join(path...)
	}
	Clean(out)
	if !src.IsDir() {
		err := os.MkdirAll(out, 0755)
		if err != nil {
			panic(err)
		}
		split := strings.Split(src.Abs(), string(os.PathSeparator))[2:]
		out = filepath.Join(append([]string{out}, split...)...)
		//out = filepath.Join(out, )
		if err = copyFileAndClose(src.FS(), 0755, src.Path(), out); err != nil {
			panic(err)
		}
	} else {
		if err := fs.WalkDir(src.FS(), ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			dstPath := filepath.Join(out, path)
			if err = copyAndClose(src.FS(), d, path, dstPath); err != nil {
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
	file, err := source.File(out)
	if err != nil {
		panic(err)
	}
	return file
}

func copyAndClose(src fs.FS, d fs.DirEntry, srcPath, dstPath string) (err error) {
	if d.IsDir() {
		return os.MkdirAll(dstPath, 0777)
	}
	info, err := d.Info()
	if err != nil {
		return
	}
	return copyFileAndClose(src, info.Mode(), srcPath, dstPath)
}

func copyFileAndClose(src fs.FS, mod fs.FileMode, srcPath, dstPath string) (err error) {
	_ = os.MkdirAll(filepath.Dir(dstPath), 0755)
	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, mod)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := src.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return
	}
	return
}

func Clean(path ...string) {
	dir := "test_data"
	if len(path) > 0 {
		dir = filepath.Join(path...)
	}
	_ = os.RemoveAll(dir)
}
