package npm

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type injector struct {
	deps []target.NodeModule
}

func Injector(deps []target.NodeModule) target.Run[target.NodeModule] {
	return injector{deps: deps}.Run
}

func (i injector) Run(ctx context.Context, m target.NodeModule, _ ...string) (err error) {
	log := plog.Get(ctx).Type(i).Set(&ctx)
	for _, module := range i.deps {
		if err = inject(log, m, module); err != nil {
			return fmt.Errorf("cannot inject dependency %s in %s: %s", module.Abs(), err, module)
		}
	}
	return
}

func inject(log plog.Logger, m target.NodeModule, lib target.NodeModule) (err error) {
	dep := lib
	nm := filepath.Join(m.Abs(), "node_modules", filepath.Base(dep.Abs()))
	log.Printf("copying module %v %v into: %s", dep.Path(), dep.Abs(), nm)
	return fs.WalkDir(dep.FS(), ".", func(s string, d fs.DirEntry, err error) error {
		log.Println("* coping file", d, s)
		if d.IsDir() {
			dst := filepath.Join(nm, s)
			if err = os.MkdirAll(dst, 0755); err != nil {
				return fmt.Errorf("os.MkdirAll: %v", err)
			}
			return nil
		}
		src, err := dep.FS().Open(s)
		if err != nil {
			return fmt.Errorf("cannot open %s: %s", s, err)
		}
		defer src.Close()
		dst, err := os.Create(filepath.Join(nm, s))
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
