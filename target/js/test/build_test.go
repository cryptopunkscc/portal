package test

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/js"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
	"testing"
)

func TestBuild(t *testing.T) {
	if err := deps.RequireBinary("npm"); err != nil {
		t.SkipNow()
	}

	ctx := context.Background()
	s := source.Dir("project")
	p, err := js.ResolveProject(s)
	test.AssertErr(t, err)

	err = npm.BuildProject().Run(ctx, p, "clean", "pack", "build")
	test.AssertErr(t, err)
}
