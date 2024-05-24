package npm

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"io"
	"io/fs"
	"os"
	"path"
)

func RunBuild(m target.NodeModule) (err error) {
	if err = exec.Run(m.Abs(), "npm", "run", "build"); err != nil {
		return fmt.Errorf("npm.RunBuild %v: %w", m.Abs(), err)
	}
	return
}

func Install(m target.NodeModule) (err error) {
	if err = exec.Run(m.Abs(), "npm", "install"); err != nil {
		return fmt.Errorf("cannot install node_modules in %s: %s", m.Abs(), err)
	}
	return
}

type Injector struct {
	deps []target.NodeModule
}

func NewInjector(deps []target.NodeModule) *Injector {
	return &Injector{deps: deps}
}

func (i Injector) Run(ctx context.Context, m target.NodeModule) (err error) {
	log := plog.Get(ctx).Type(i).Set(&ctx)
	for _, module := range i.deps {
		if err = inject(log, m, module); err != nil {
			return fmt.Errorf("cannot inject dependency %s in %s: %s", module.Abs(), err, module)
		}
	}
	return
}

func inject(log plog.Logger, m target.NodeModule, lib target.NodeModule) (err error) {
	dep := lib.Lift()
	nm := path.Join(m.Abs(), "node_modules", path.Base(dep.Abs()))
	log.Printf("copying module %v %v into: %s", dep.Path(), dep.Abs(), nm)
	return fs.WalkDir(dep.Files(), ".", func(s string, d fs.DirEntry, err error) error {
		log.Println("* coping file", d, s)
		if d.IsDir() {
			dst := path.Join(nm, s)
			if err = os.MkdirAll(dst, 0755); err != nil {
				return fmt.Errorf("os.MkdirAll: %v", err)
			}
			return nil
		}
		src, err := dep.Files().Open(s)
		if err != nil {
			return fmt.Errorf("cannot open %s: %s", s, err)
		}
		defer src.Close()
		dst, err := os.Create(path.Join(nm, s))
		if err != nil {
			return fmt.Errorf("os.Create: %v", err)
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			return fmt.Errorf("io.Copy: %v", err)
		}
		return nil
	})
}
