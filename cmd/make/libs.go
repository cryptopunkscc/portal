package main

import (
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/target"
	npm2 "github.com/cryptopunkscc/portal/target2/npm"
	"github.com/cryptopunkscc/portal/target2/source"
	"log"
)

func (d *Install) buildJsLibs() {
	libs, err := source.File(d.root, "target", "js")
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range target.List(npm2.Resolve, libs) {
		if !p.PkgJson().CanBuild() {
			continue
		}
		if err := npm.Install(p); err != nil {
			log.Fatalln(err)
		}
		if err := npm.RunBuild(p); err != nil {
			log.Fatalln(err)
		}
	}
}
