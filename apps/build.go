package apps

import (
	"context"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/target/all"
	"path/filepath"
)

func Build(args ...string) error {
	ctx := context.Background()
	root, err := golang.FindProjectRoot()
	if err != nil {
		return err
	}
	appsDir := filepath.Join(root, "apps")
	return all.BuildRecursive(ctx, appsDir, args...)
}
