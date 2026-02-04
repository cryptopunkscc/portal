package env

import (
	"os"
	"path/filepath"
)

type Key string

func (k Key) Default(get func() string) {
	if !k.Exist() {
		v := get()
		k.Set(v)
	}
}

func (k Key) Exist() bool {
	return len(k.Get()) > 0
}

func (k Key) Get() string {
	return os.Getenv(string(k))
}

func (k Key) Set(v string) {
	err := os.Setenv(string(k), v)
	if err != nil {
		panic(err)
	}
}

func (k Key) Unset() {
	if err := os.Unsetenv(string(k)); err != nil {
		panic(err)
	}
}

func (k Key) SetDir(dir string, path ...string) {
	dir = filepath.Join(dir, filepath.Join(path...))
	k.Set(dir)
}

func (k Key) MkdirAll() (dir string) {
	abs := k.Get()
	if abs == "" {
		panic("environment variable " + k + " not set")
	}
	err := os.MkdirAll(abs, 0755)
	if err != nil {
		panic(err)
	}
	return abs
}
