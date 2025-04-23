package js_lib_build

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	npm2 "github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/npm"
	"log"
	"os"
)

func Run() (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot resolve working dir: %v", err)
	}

	libs, err := source.File(wd, "core", "js")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	for _, p := range target.Any[target.NodeModule](
		target.Skip("node_modules"),
		target.Try(npm2.ResolveNodeModule),
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
	return
}
