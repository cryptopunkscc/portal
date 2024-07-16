package dist

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/target"
	"os"
	"path"
)

func Dist(ctx context.Context, project target.Project) (err error) {
	if err = CopyIcon(ctx, project); err != nil {
		return
	}
	if err = CopyManifest(ctx, project); err != nil {
		return
	}
	return
}

func CopyIcon(_ context.Context, project target.Project) (err error) {
	if project.Manifest().Icon == "" {
		return
	}
	iconSrc := path.Join(project.Abs(), project.Manifest().Icon)
	iconName := "icon" + path.Ext(project.Manifest().Icon)
	iconDst := path.Join(project.Abs(), "dist", iconName)
	if err = fs2.CopyFile(iconSrc, iconDst); err != nil {
		return
	}
	project.Manifest().Icon = iconName
	return
}

func CopyManifest(_ context.Context, project target.Project) (err error) {
	bytes, err := json.Marshal(project.Manifest())
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(project.Abs(), "dist", target.PortalJsonFilename), bytes, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %v", err)
	}
	return
}
