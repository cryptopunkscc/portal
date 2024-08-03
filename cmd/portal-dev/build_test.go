package main

import (
	"context"
	. "github.com/cryptopunkscc/portal/target"
	"testing"
)

func TestBuild(t *testing.T) {
	mod := Module[Portal_]{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := mod.FeatBuild().Run(ctx, "../../apps")
	if err != nil {
		t.Fatal(err)
	}
}
