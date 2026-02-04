package os

import (
	path2 "path"
	"path/filepath"
)

func Abs(path ...string) string {
	src := path2.Join(path...)
	abs, err := filepath.Abs(src)
	if err != nil {
		return src
	}
	return abs
}
