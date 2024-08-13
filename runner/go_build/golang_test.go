package go_build

import (
	"context"
	golang "github.com/cryptopunkscc/portal/resolve/go"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoRunner_Run(t *testing.T) {
	ctx := context.Background()
	file, err := source.File("test", "test")
	if err != nil {
		t.Fatal(err)
	}
	project, err := golang.ResolveProject(file)
	if err != nil {
		t.Fatal(err)
	}
	expected := target.Builds{
		"default": target.Build{Deps: []string(nil), Env: []string(nil), Cmd: "go build -o dist/main"},
		"windows": target.Build{Deps: []string(nil), Env: []string(nil), Cmd: "go build -o dist/main.exe"},
	}
	assert.Equal(t, expected, project.Build())
	run := Runner()
	if err = run(ctx, project); err != nil {
		t.Fatal(err)
	}
}
