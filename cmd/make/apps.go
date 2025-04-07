package main

import (
	"context"
	"github.com/cryptopunkscc/portal/runner/any_build"
	"log"
	"path/filepath"
)

func (d *Install) buildEmbedApps(platforms ...string) {
	appsDir := filepath.Join(d.root, "apps")
	if err := any_build.Run(context.TODO(), appsDir); err != nil {
		log.Fatal(err)
	}
}
