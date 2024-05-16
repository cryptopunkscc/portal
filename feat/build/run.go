package build

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	js "github.com/cryptopunkscc/go-astral-js/pkg/binding/out"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/dist"
	"github.com/cryptopunkscc/go-astral-js/runner/pack"
	"path"
)

func Run(dir string) (err error) {
	libs := array.FromChan(project.FindInFS[target.NodeModule](js.PortalLibFS))
	if err = Dist(dir, ".", libs...); err != nil {
		return fmt.Errorf("cannot build portal apps: %w", err)
	}
	if err = Pack(dir, "."); err != nil {
		return fmt.Errorf("cannot bundle portal apps: %w", err)
	}
	return
}

func Dist(root, dir string, dependencies ...target.NodeModule) (err error) {
	for m := range project.FindInPath[target.Project](path.Join(root, dir)) {
		if !m.CanNpmRunBuild() {
			continue
		}
		if err = dist.Run(m, dependencies); err != nil {
			return err
		}
	}
	return
}

func Pack(base, sub string) (err error) {
	err = errors.New("no targets found")
	for app := range project.FindInPath[*project.PortalRawModule](path.Join(base, sub)) {
		if err = pack.Run(app); err != nil {
			return fmt.Errorf("bundle target %v: %v", app.Path(), err)
		}
	}
	return
}
