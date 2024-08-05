package main

import (
	"context"
	"github.com/cryptopunkscc/portal/feat/build"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/clean"
	"github.com/cryptopunkscc/portal/runner/go_build"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/npm_build"
	"github.com/cryptopunkscc/portal/runner/pack"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"path/filepath"
)

func (d *Install) buildEmbedApps(platforms ...string) {
	file, err := source.File(d.root, "target", "js", "embed", "portal")
	if err != nil {
		log.Fatal(err)
	}
	jsLibs := target.List(target.Any[target.NodeModule](
		target.Skip("node_modules"),
		target.Try(npm.Resolve),
	), file)

	feat := build.NewFeat(
		clean.NewRunner().Call,
		multi.NewRunner[target.Project_](
			go_build.NewRun(platforms...).Portal(),
			npm_build.NewRun(jsLibs...).Portal(),
		).Run,
		pack.Run,
	)
	appsDir := filepath.Join(d.root, "apps")
	if err := feat.Run(context.TODO(), appsDir); err != nil {
		log.Fatal(err)
	}
}
