package apps_build

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/build"
	"github.com/cryptopunkscc/portal/runner/clean"
	"github.com/cryptopunkscc/portal/runner/go_build"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/npm_build"
	"github.com/cryptopunkscc/portal/runner/pack"
	"os"
	"path/filepath"
)

func Run() error {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := source.File(wd, "core", "js", "embed", "portal")
	if err != nil {
		return err
	}
	jsLibs := target.Any[target.NodeModule](
		target.Skip("node_modules"),
		target.Try(npm.Resolve),
	).List(file)

	feat := build.NewRunner(
		clean.Runner(),
		multi.NewRun[target.Project_](
			go_build.Runner( /*TODO*/ ).Portal(),
			npm_build.Runner(jsLibs...).Portal(),
		),
		pack.Run,
	)
	appsDir := filepath.Join(wd, "apps")
	if err = feat.Run(ctx, appsDir); err != nil {
		return err
	}
	return nil
}
