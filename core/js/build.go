package js

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
	"log"
)

func BuildPortalLib() (err error) {
	defer plog.TraceErr(&err)

	wd, err := golang.FindProjectRoot()
	if err != nil {
		return
	}

	libs, err := source.File(wd, "core", "js")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	for _, p := range target.Any[target.NodeModule](
		target.Skip("node_modules"),
		target.Try(npm.ResolveNodeModule),
	).List(libs) {
		if !p.PkgJson().CanBuild() {
			continue
		}
		log.Printf("building js libs for %s", p.Abs())
		if err = npm.Install(ctx, p); err != nil {
			return
		}
		if err = npm.BuildModule(ctx, p); err != nil {
			return
		}
	}
	return
}
