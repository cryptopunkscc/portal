package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"os"
	"path/filepath"
)

func Dist(ctx context.Context, project target.Project_) (err error) {
	if err = copyIcon(ctx, project); err != nil {
		return
	}
	if err = copyManifest(ctx, project); err != nil {
		return
	}
	return
}

func copyIcon(_ context.Context, project target.Project_) (err error) {
	if project.Manifest().Icon == "" {
		return
	}
	iconSrc := filepath.Join(project.Abs(), project.Manifest().Icon)
	iconName := "icon" + filepath.Ext(project.Manifest().Icon)
	iconDst := filepath.Join(project.Abs(), "dist", iconName)
	if err = fs2.CopyFile(iconSrc, iconDst); err != nil {
		return
	}
	project.Manifest().Icon = iconName
	return
}

func copyManifest(_ context.Context, project target.Project_) (err error) {
	bytes, err := json.Marshal(project.Manifest())
	if err != nil {
		return err
	}
	name := filepath.Join(project.Abs(), "dist", manifest.AppFilename+".json")
	if err = os.WriteFile(name, bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}

func Dist2(ctx context.Context, project target.Project_, platform ...string) (err error) {
	if err = copyIcon(ctx, project); err != nil {
		return
	}
	distManifest := buildDistManifest(project, platform...)
	if err = writeDistManifest(project, distManifest); err != nil {
		return
	}
	return
}

func buildDistManifest(project target.Project_, platform ...string) (out manifest.Dist) {
	b := project.Build().Get(platform...)
	out = manifest.Dist{
		App:    *project.Manifest(),
		Api:    *project.Api(),
		Config: *project.Config(),
		Target: b.Target,
		Release: manifest.Release{
			Version: 0, // TODO
		},
	}
	return
}

func writeDistManifest(project target.Project_, distManifest manifest.Dist) (err error) {
	bytes, err := json.Marshal(distManifest)
	if err != nil {
		return err
	}
	name := filepath.Join(project.Abs(), "dist", manifest.AppFilename+".json")
	if err = os.WriteFile(name, bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
