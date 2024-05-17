package dist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/runner/npm"
	"github.com/cryptopunkscc/go-astral-js/target"
	"os"
	"path"
)

type Runner struct {
	dependencies []target.NodeModule
}

func NewRunner(dependencies []target.NodeModule) *Runner {
	return &Runner{dependencies: dependencies}
}

func (r Runner) Run(ctx context.Context, m target.Project) (err error) {
	if err = r.Prepare(ctx, m); err != nil {
		return fmt.Errorf("dist.Prepare: %w", err)
	}
	if err = r.Dist(ctx, m); err != nil {
		return fmt.Errorf("dist.Dist: %w", err)
	}
	return
}

func (r Runner) Prepare(ctx context.Context, m target.Project) (err error) {
	if err = npm.Install(m); err != nil {
		return
	}
	if err = npm.NewInjector(r.dependencies).Run(ctx, m); err != nil {
		return
	}
	return
}

func (r Runner) Dist(ctx context.Context, m target.Project) (err error) {
	if !m.PkgJson().CanBuild() {
		return errors.New("missing npm build in package.json")
	}
	if err = npm.RunBuild(m); err != nil {
		return
	}
	if err = r.CopyIcon(ctx, m); err != nil {
		return
	}
	if err = r.CopyManifest(ctx, m); err != nil {
		return
	}
	return
}

func (r Runner) CopyIcon(_ context.Context, m target.Project) (err error) {
	if m.Manifest().Icon == "" {
		return
	}
	iconSrc := path.Join(m.Abs(), m.Manifest().Icon)
	iconName := "icon" + path.Ext(m.Manifest().Icon)
	iconDst := path.Join(m.Abs(), "dist", iconName)
	if err = fs.CopyFile(iconSrc, iconDst); err != nil {
		return
	}
	m.Manifest().Icon = iconName
	return
}

func (r Runner) CopyManifest(_ context.Context, m target.Project) (err error) {
	bytes, err := json.Marshal(m.Manifest())
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(m.Abs(), "dist", target.PortalJsonFilename), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
