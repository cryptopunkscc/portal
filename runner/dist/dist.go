package dist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"os"
	"path/filepath"
)

func Dist(ctx context.Context, project target.Project_) (err error) {
	if err = CopyIcon(ctx, project); err != nil {
		return
	}
	if err = CopyManifest(ctx, project); err != nil {
		return
	}
	return
}

func CopyIcon(_ context.Context, project target.Project_) (err error) {
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

func CopyManifest(_ context.Context, project target.Portal_) (err error) {
	bytes, err := json.Marshal(project.Manifest())
	if err != nil {
		return err
	}
	name := filepath.Join(project.Abs(), "dist", target.ManifestFilename+".json")
	if err = os.WriteFile(name, bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
