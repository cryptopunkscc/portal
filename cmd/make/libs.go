package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	npm2 "github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/npm"
	"log"
)

func (d *Install) buildJsLibs() {
	ctx := context.Background()
	libs, err := source.File(d.root, "runtime", "js")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range target.Any[target.NodeModule](
		target.Skip("node_modules"),
		target.Try(npm2.Resolve),
	).List(libs) {
		if !p.PkgJson().CanBuild() {
			continue
		}
		log.Printf("building js libs for %s", p.Abs())
		if err := npm.Install(ctx, p); err != nil {
			log.Fatalln(err)
		}
		if err := npm.Build(ctx, p); err != nil {
			log.Fatalln(err)
		}
	}
}
