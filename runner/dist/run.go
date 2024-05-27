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

func NewRun(dependencies []target.NodeModule) target.Run[target.Project] {
	return NewRunner(dependencies).Run
}

type Runner struct {
	dependencies []target.NodeModule
}

func NewRunner(dependencies []target.NodeModule) *Runner {
	return &Runner{dependencies: dependencies}
}

func (r Runner) Run(ctx context.Context, project target.Project) (err error) {
	if err = r.Prepare(ctx, project); err != nil {
		return fmt.Errorf("dist.Prepare: %w", err)
	}
	if err = r.Dist(ctx, project); err != nil {
		return fmt.Errorf("dist.Dist: %w", err)
	}
	return
}

func (r Runner) Prepare(ctx context.Context, project target.Project) (err error) {
	if err = npm.Install(project); err != nil {
		return
	}
	if err = npm.NewInjector(r.dependencies).Run(ctx, project); err != nil {
		return
	}
	return
}

func (r Runner) Dist(ctx context.Context, project target.Project) (err error) {
	if !project.PkgJson().CanBuild() {
		return errors.New("missing npm build in package.json")
	}
	if err = npm.RunBuild(project); err != nil {
		return
	}
	if err = r.CopyIcon(ctx, project); err != nil {
		return
	}
	if err = r.CopyManifest(ctx, project); err != nil {
		return
	}
	return
}

func (r Runner) CopyIcon(_ context.Context, project target.Project) (err error) {
	if project.Manifest().Icon == "" {
		return
	}
	iconSrc := path.Join(project.Abs(), project.Manifest().Icon)
	iconName := "icon" + path.Ext(project.Manifest().Icon)
	iconDst := path.Join(project.Abs(), "dist", iconName)
	if err = fs.CopyFile(iconSrc, iconDst); err != nil {
		return
	}
	project.Manifest().Icon = iconName
	return
}

func (r Runner) CopyManifest(_ context.Context, project target.Project) (err error) {
	bytes, err := json.Marshal(project.Manifest())
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(project.Abs(), "dist", target.PortalJsonFilename), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
