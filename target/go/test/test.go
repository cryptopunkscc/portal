package test

import (
	"context"
	"embed"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/test"
	golang "github.com/cryptopunkscc/portal/target/go"
	"github.com/cryptopunkscc/portal/target/source"
	"testing"
)

//go:embed project
var goProjectFS embed.FS

var dir = source.Dir(".")

func ResolveGoProject(t *testing.T, path ...string) {
	if len(path) == 0 {
		path = append(path, ".")
	}
	dir := source.Dir(path...)
	project, err := golang.ResolveProject(dir)
	test.AssertErr(t, err)
	println(target.Sprint(project))
}

func BuildGoProject(t *testing.T) {
	ctx := context.Background()
	project, err := golang.ResolveProject(dir)
	test.AssertErr(t, err)
	err = golang.BuildProject().Run(ctx, project, "clean", "pack")
	test.AssertErr(t, err)
}
