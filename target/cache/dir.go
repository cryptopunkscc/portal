package cache

import (
	"os"
	"path"
)

type Deps interface {
	Executable() string
}

func Dir(deps Deps) (dir string) {
	var err error
	if dir, err = os.UserCacheDir(); err != nil {
		panic(err)
	}
	dir = path.Join(dir, deps.Executable())
	return
}
