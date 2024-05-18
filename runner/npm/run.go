package npm

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io"
	"io/fs"
	"log"
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

func InjectDependencies(m target.NodeModule, deps []target.NodeModule) (err error) {
	for _, module := range deps {
		if err = InjectDependency(m, module); err != nil {
			return fmt.Errorf("cannot inject dependency %s in %s: %s", module.Abs(), err, module)
		}
	}
	return
}

func InjectDependency(m target.NodeModule, dep target.NodeModule) (err error) {
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
