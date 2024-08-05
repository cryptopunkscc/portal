package clean

import (
	"context"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var Defaults = []string{"build", "dist", "node_modules"}

type Runner struct{ names map[string]any }

func NewRunner(names ...string) (r *Runner) {
	if len(names) == 0 {
		names = Defaults
	}
	r = &Runner{names: make(map[string]any)}
	for _, name := range names {
		r.names[name] = name
	}
	return
}

func (r Runner) match(name string) bool { return r.names[name] != nil }

func (r Runner) Run(_ context.Context, src target.Source) error { return r.Call(src.Abs()) }
func (r Runner) Call(dir string) error {
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && r.match(info.Name()) {
			if err = os.RemoveAll(path); err != nil {
				return err
			}
			log.Println("* removed:", path)
			return fs.SkipDir
		}
		return nil
	})
}
