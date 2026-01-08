package os

import (
	"os"
	path2 "path"
	"path/filepath"
)

func Abs(path ...string) string {
	src := path2.Join(path...)
	if path2.IsAbs(src) {
		return src
	}
	base, err := os.Getwd()
	if err != nil {
		return src
	}
	base = filepath.ToSlash(base)
	return path2.Join(base, src)
}
