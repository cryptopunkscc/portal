package target

import (
	"os"
	"path/filepath"
)

func Abs(path ...string) string {
	src := filepath.Join(path...)
	if filepath.IsAbs(src) {
		return src
	}
	base, err := os.Getwd()
	if err != nil {
		return src
	}
	return filepath.Join(base, src)
}

func CacheDir(name string) (dir string) {
	var err error
	if dir, err = os.UserCacheDir(); err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, name)
	return
}
