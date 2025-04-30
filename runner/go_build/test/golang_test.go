package test

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	golang "github.com/cryptopunkscc/portal/resolve/go"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/go_build"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
)

func TestGoRunner_Run(t *testing.T) {
	defer test.Clean()
	sub, err := fs.Sub(goProject, "go")
	test.AssertErr(t, err)
	s := source.Embed(sub)
	src := test.Copy(s)
	ctx := context.Background()
	project, err := golang.ResolveProject(src)
	if err != nil {
		t.Fatal(err)
	}

	build := target.Builds{
		"default": target.Build{Out: "main", Cmd: "go build -o dist/main"},
		"linux":   target.Build{Out: "main", Cmd: "go build -o dist/main", Deps: []string{"gcc", "libgtk-3-dev", "libayatana-appindicator3-dev"}},
		"windows": target.Build{Out: "main.exe", Cmd: "go build -ldflags -H=windowsgui -o dist/main.exe", Env: []string{"CGO_ENABLED=1"}},
	}
	assert.Equal(t, build, project.Build())
	run := go_build.NewRun()
	if err = run(ctx, project); err != nil {
		t.Fatal(err)
	}
}
