package clean

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var Defaults = []string{"build", "dist", "node_modules"}

type runner struct{ names map[string]any }

func Runner(names ...string) func(string) error {
	if len(names) == 0 {
		names = Defaults
	}
	r := &runner{names: make(map[string]any)}
	for _, name := range names {
		r.names[name] = name
	}
	return r.call
}

func (r runner) match(name string) bool { return r.names[name] != nil }

func (r runner) call(dir string) error {
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
