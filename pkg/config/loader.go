package config

import (
	"errors"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
	"path/filepath"
)

type Loader[T any] struct {
	Unmarshal func([]byte, any) error
	Config    T
	wd        string
	File      string
	Dir       string
}

func (l *Loader[T]) Load(path ...string) (err error) {
	defer plog.TraceErr(&err)
	p := ""
	if p, err = l.abs(path...); err != nil {
		return
	}
	isDirectory := false
	if isDirectory, err = isDir(p); err != nil {
		return
	}
	if !isDirectory {
		if err = l.load(p); err != nil {
			return
		}
		l.File = filepath.Base(p)
		l.Dir = filepath.Dir(p)
		return
	}
	for {
		l.Dir = p
		if err = l.load(l.Dir, l.File); err == nil {
			return
		}
		if p = filepath.Dir(l.Dir); p == l.Dir {
			return ErrNotFound // end search when the root directory is reached
		}
	}
}

var ErrNotFound = errors.New("config not found")

func (l *Loader[T]) load(path ...string) (err error) {
	defer plog.TraceErr(&err)
	var bytes []byte
	if bytes, err = os.ReadFile(filepath.Join(path...)); err != nil {
		return
	}
	if err = l.Unmarshal(bytes, l.Config); err != nil {
		return
	}
	return
}

func (l *Loader[T]) abs(path ...string) (p string, err error) {
	p = filepath.Join(path...)
	if !filepath.IsAbs(p) {
		if len(l.wd) == 0 {
			if l.wd, err = os.Getwd(); err != nil {
				return
			}
		}
		p = filepath.Join(l.wd, p)
	}
	return
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, plog.Err(err)
	}
	return info.IsDir(), nil
}
