package dist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/fs"
	"github.com/cryptopunkscc/portal/target"
	"os"
	"path"
)

func NewRun(dependencies []target.NodeModule) target.Run[target.Project] {
	return Runner{
		NpmRunner: NewNpmRunner(dependencies),
		GoRunner:  NewGoRunner(),
	}.Run
}

type Runner struct {
	NpmRunner
	GoRunner
}

func (r Runner) Run(ctx context.Context, project target.Project) (err error) {
	// TODO replace switch with injected factory
	switch v := project.(type) {
	case target.ProjectNpm:
		if err = r.NpmRunner.Run(ctx, v); err != nil {
			return
		}
	case target.ProjectGo:
		if err = r.GoRunner.Run(ctx, v); err != nil {
			return
		}
	}
	if err = r.Dist(ctx, project); err != nil {
		return
	}
	return
}

func (r Runner) Dist(ctx context.Context, project target.Project) (err error) {
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
