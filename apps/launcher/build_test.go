package launcher

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
	"testing"
)

func TestBuild(t *testing.T) {
	t.SkipNow()
	ctx := context.Background()
	dir := source.Dir("svelte")
	project, err := npm.Resolve_(dir)
	test.AssertErr(t, err)
	err = npm.BuildProject().Run(ctx, project, "clean", "pack")
	test.AssertErr(t, err)
}
