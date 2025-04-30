package apps

import (
	"context"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/runner/any_build"
	"path/filepath"
)

func Build(args ...string) error {
	ctx := context.Background()
	root, err := golang.FindProjectRoot()
	if err != nil {
		return err
	}
	appsDir := filepath.Join(root, "apps")
	return any_build.Run(ctx, appsDir, args...)
}
