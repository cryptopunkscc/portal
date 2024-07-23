package cache

import (
	"os"
	"path/filepath"
)

type Deps interface {
	Executable() string
}

func Dir(deps Deps) (dir string) {
	var err error
	if dir, err = os.UserCacheDir(); err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, deps.Executable())
	return
}
