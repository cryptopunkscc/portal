package build

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	"github.com/cryptopunkscc/go-astral-js/runner/dist"
	"github.com/cryptopunkscc/go-astral-js/runner/pack"
	"github.com/cryptopunkscc/go-astral-js/target"
	js "github.com/cryptopunkscc/go-astral-js/target/js/embed"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
	"path"
)

type Feat struct {
	dependencies []target.NodeModule
}

func NewFeat(dependencies ...target.NodeModule) *Feat {
	if len(dependencies) == 0 {
		dependencies = array.FromChan(sources.FromFS[target.NodeModule](js.PortalLibFS))
	}
	return &Feat{dependencies: dependencies}
}

func (r Feat) Run(ctx context.Context, dir string) (err error) {
	if err = r.Dist(ctx, dir); err != nil {
		return fmt.Errorf("cannot build portal apps: %w", err)
	}
	if err = r.Pack(ctx, dir, "."); err != nil {
		return fmt.Errorf("cannot bundle portal apps: %w", err)
	}
	return
}

func (r Feat) Dist(ctx context.Context, dir ...string) (err error) {
	for m := range sources.FromPath[target.Project](path.Join(dir...)) {
		if !m.PkgJson().CanBuild() {
			continue
		}
		if err = dist.NewRunner(r.dependencies).Run(ctx, m); err != nil {
			return fmt.Errorf("build.Dist: %w", err)
		}
	}
	return
}

func (r Feat) Pack(ctx context.Context, base, sub string) (err error) {
	err = errors.New("no targets found")
	for app := range sources.FromPath[target.Dist](path.Join(base, sub)) {
		if err = pack.Run(ctx, app); err != nil {
			return fmt.Errorf("bundle target %v: %v", app.Path(), err)
		}
	}
	return
}
