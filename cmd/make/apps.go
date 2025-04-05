package main

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
	"log"
	"path/filepath"
)

func (d *Install) buildEmbedApps(platforms ...string) {
	file, err := source.File(d.root, "core", "js", "embed", "portal")
	if err != nil {
		log.Fatal(err)
	}
	jsLibs := target.Any[target.NodeModule](
		target.Skip("node_modules"),
		target.Try(npm.Resolve),
	).List(file)

	feat := build.NewRunner(
		clean.Runner(),
		multi.NewRun[target.Project_](
			go_build.Runner(platforms...).Portal(),
			npm_build.Runner(jsLibs...).Portal(),
		),
		pack.Run,
	)
	appsDir := filepath.Join(d.root, "apps")
	if err := feat.Run(context.TODO(), appsDir); err != nil {
		log.Fatal(err)
	}
}
