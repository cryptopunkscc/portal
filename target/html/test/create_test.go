package test

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/html"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateProject(t *testing.T) {
	name := "created_project"
	_ = os.RemoveAll(name)

	opts := npm.CreateOpts{}
	opts.Template = "svelte"
	opts.Path = name

	err := html.Create(opts)
	test.AssertErr(t, err)

	p, err := html.ResolveProject.Resolve(source.Dir(name))
	test.AssertErr(t, err)
	assert.NotNil(t, p)

	ctx := context.Background()
	err = npm.BuildProject().Run(ctx, p)
	plog.Println(err)
}

func TestCreateDist(t *testing.T) {
	name := "created_dist"
	_ = os.RemoveAll(name)

	opts := npm.CreateOpts{}
	opts.Template = "html"
	opts.Path = name
	err := html.Create(opts)
	test.AssertErr(t, err)

	_, err = html.ResolveDist.Resolve(source.Dir(name))
	test.AssertErr(t, err)
}
