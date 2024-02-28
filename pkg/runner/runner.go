package runner

import (
	"io/fs"
	"log"
	"os"
	"path"
)

type Runner struct {
	Frontends []Target
	Backends  []Target
}

func New(dir string, resolve GetTargets) (out *Runner, err error) {
	var targets []Target
	targets, err = resolve(dir)
	if err != nil {
		return
	}
	out = &Runner{}
	for _, d := range targets {
		switch {
		case isBackend(d.Files):
			log.Println("found backend:", d.Path)
			out.Backends = append(out.Backends, d)
		case isFrontend(d.Files):
			log.Println("found frontend:", d.Path)
			out.Frontends = append(out.Frontends, d)
		}
	}
	return
}

func isFrontend(dir fs.FS) bool {
	stat, err := fs.Stat(dir, "index.html")
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func isBackend(dir fs.FS) bool {
	_, err := fs.Stat(dir, "index.html")
	return err != nil
}

func ResolveSrc(dir string, name string) (f string, err error) {
	f = path.Join(dir, "dist", name)
	if _, err = os.Stat(f); err == nil {
		return
	}
	f = path.Join(dir, name)
	if _, err = os.Stat(f); err == nil {
		return
	}
	return
}
