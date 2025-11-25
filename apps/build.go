package apps

import (
	"context"
	"path/filepath"

	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/target/all"
)

func Build(args ...string) error {
	ctx := context.Background()
	appsDir, err := Dir()
	if err != nil {
		return err
	}
	args = append(args, appsDir)
	return all.BuildRecursive(ctx, appsDir, args...)
}

func Dir() (d string, err error) {
	root, err := golang.FindProjectRoot()
	if err != nil {
		return
	}
	d = filepath.Join(root, "apps")
	return
}
