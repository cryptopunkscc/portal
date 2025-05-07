package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/pkg/plog"
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

func Dist2(ctx context.Context, project target.Project_, target manifest.Target) (err error) {
	defer plog.TraceErr(&err)
	if err = copyIcon(ctx, project); err != nil {
		return
	}
	dist := buildDistManifest(project, target)
	if err = writeDistManifest(project, dist); err != nil {
		return
	}
	return
}

func buildDistManifest(project target.Project_, target manifest.Target) (out manifest.Dist) {
	out = manifest.Dist{
		App:    *project.Manifest(),
		Api:    *project.Api(),
		Config: *project.Config(),
		Target: target,
		Release: manifest.Release{
			Version: 0, // TODO
		},
	}
	return
}

func writeDistManifest(project target.Project_, dist manifest.Dist) (err error) {
	defer plog.TraceErr(&err)

	bytes, err := json.Marshal(dist)
	if err != nil {
		return err
	}

	path := []string{project.Abs(), "dist"}
	path = append(path, DistPath(dist.Target)...)
	path = append(path, manifest.AppFilename+".json")
	name := filepath.Join(path...)

	return os.WriteFile(name, bytes, 0644)
}

func DistPath(target manifest.Target) (path []string) {
	if len(target.OS) > 0 {
		path = append(path, target.OS)
	}
	if len(target.Arch) > 0 {
		path = append(path, target.Arch)
	}
	return
}
