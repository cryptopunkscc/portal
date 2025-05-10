package test

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/test"
	golang "github.com/cryptopunkscc/portal/target/go"
	"github.com/cryptopunkscc/portal/target/source"
	"testing"
)

func TestBuildProject(t *testing.T) {
	d, err := source.Embed(goProjectFS).Sub("project")
	test.AssertErr(t, err)

	n := test.CleanDir(t, ".test_project")
	err = d.CopyTo(n)
	test.AssertErr(t, err)

	d = source.Dir(n)
	project, err := golang.ResolveProject(d)
	test.AssertErr(t, err)

	ctx := context.Background()
	err = golang.BuildProject("linux", "linux/arm64", "windows").Run(ctx, project, "clean", "pack")
	test.AssertErr(t, err)
}
