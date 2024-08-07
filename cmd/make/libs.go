package main

import (
	npm2 "github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/target"
	"log"
)

func (d *Install) buildJsLibs() {
	libs, err := source.File(d.root, "target", "js")
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
		if err := npm.Install(p); err != nil {
			log.Fatalln(err)
		}
		if err := npm.RunBuild(p); err != nil {
			log.Fatalln(err)
		}
	}
}
