package build

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/array"
	js "github.com/cryptopunkscc/go-astral-js/pkg/js/embed"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/dist"
	"github.com/cryptopunkscc/go-astral-js/runner/pack"
	"path"
)

type Feat struct {
	dependencies []target.NodeModule
}

func NewFeat(dependencies ...target.NodeModule) *Feat {
	if len(dependencies) == 0 {
		dependencies = array.FromChan(project.FromFS[target.NodeModule](js.PortalLibFS))
	}
	return &Feat{dependencies: dependencies}
}

func (r Feat) Run(ctx context.Context, dir string) (err error) {
	if err = r.Dist(ctx, dir, "."); err != nil {
		return fmt.Errorf("cannot build portal apps: %w", err)
	}
	if err = r.Pack(ctx, dir, "."); err != nil {
		return fmt.Errorf("cannot bundle portal apps: %w", err)
	}
	return
}

func (r Feat) Dist(ctx context.Context, root, dir string) (err error) {
	for m := range project.FromPath[target.Project](path.Join(root, dir)) {
		if !m.CanNpmRunBuild() {
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
	for app := range project.FromPath[target.Dist](path.Join(base, sub)) {
		if err = pack.Run(ctx, app); err != nil {
			return fmt.Errorf("bundle target %v: %v", app.Path(), err)
		}
	}
	return
}
