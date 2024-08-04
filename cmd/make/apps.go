package main

import (
	"context"
	"github.com/cryptopunkscc/portal/feat/build"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/all_build"
	"github.com/cryptopunkscc/portal/runner/pack"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"path/filepath"
)

func (d *Install) buildEmbedApps(platforms ...string) {
	buildEmbedApps[target.Project_](d.root, platforms...)
}

func buildEmbedApps[T target.Portal_](root string, platforms ...string) {
	file, err := source.File(root, "target", "js", "embed", "portal")
	if err != nil {
		log.Fatal(err)
	}
	feat := build.NewFeat(
		sources.Resolver[T](),
		all_build.NewRun,
		pack.Run,
		target.List(npm.Resolve, file),
		platforms...,
	)
	appsDir := filepath.Join(root, "apps")
	if err := feat.Run(context.TODO(), appsDir); err != nil {
		log.Fatal(err)
	}
}
