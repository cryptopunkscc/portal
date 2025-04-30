package test

import (
	"github.com/cryptopunkscc/portal/api/target"
	golang "github.com/cryptopunkscc/portal/target/go"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolve(t *testing.T) {
	s, err := source.Embed(goProjectFS).Sub("project")
	test.AssertErr(t, err)

	p, err := golang.ResolveProject(s)
	test.AssertErr(t, err)

	build := target.Builds{
		"default": target.Build{Out: "main", Cmd: "go build -o dist/main"},
		"linux":   target.Build{Out: "main", Cmd: "go build -o dist/main", Deps: []string{"gcc", "libgtk-3-dev", "libayatana-appindicator3-dev"}},
		"windows": target.Build{Out: "main.exe", Cmd: "go build -ldflags -H=windowsgui -o dist/main.exe", Env: []string{"CGO_ENABLED=1"}},
	}
	assert.Equal(t, build, p.Build())
	assert.Equal(t, "go", p.Manifest().Schema)
}
