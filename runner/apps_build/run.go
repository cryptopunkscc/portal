package apps_build

import (
	"context"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/runner/any_build"
	"os"
	"path/filepath"
)

func Run(args ...string) error {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	root, err := golang.FindProjectRoot(wd)
	if err != nil {
		return err
	}
	appsDir := filepath.Join(root, "apps")
	return any_build.Run(ctx, appsDir, args...)
}
