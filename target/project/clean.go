package project

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var Defaults = []string{"build", "dist", "node_modules"}

func Cleaner(names ...string) func(string) error {
	if len(names) == 0 {
		names = Defaults
	}
	r := &cleaner{names: make(map[string]any)}
	for _, name := range names {
		r.names[name] = name
	}
	return r.call
}

type cleaner struct{ names map[string]any }

func (r cleaner) match(name string) bool { return r.names[name] != nil }

func (r cleaner) call(dir string) error {
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && r.match(info.Name()) {
			if err = os.RemoveAll(path); err != nil {
				return err
			}
			log.Println("* clean:", path)
			return fs.SkipDir
		}
		return nil
	})
}
