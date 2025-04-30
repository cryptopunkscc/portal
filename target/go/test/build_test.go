package test

import (
	"context"
	golang "github.com/cryptopunkscc/portal/target/go"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"testing"
)

func TestBuildProject(t *testing.T) {

	d, err := source.Embed(goProjectFS).Sub("project")
	test.AssertErr(t, err)

	project, err := golang.ResolveProject(d)
	test.AssertErr(t, err)

	ctx := context.Background()
	err = golang.BuildProject().Run(ctx, project)
	test.AssertErr(t, err)
}
