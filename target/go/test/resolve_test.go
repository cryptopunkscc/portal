package test

import (
	"testing"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/test"
	golang "github.com/cryptopunkscc/portal/target/go"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	s, err := source.Embed(goProjectFS).Sub("project")
	test.AssertErr(t, err)

	p, err := golang.ResolveProject(s)
	test.AssertErr(t, err)

	expected := manifest.Builds{
		Builds: manifest.InnerBuilds{
			Default: manifest.Build{
				Out:  "main",
				Cmd:  "go build -o $OUT",
				Exec: "go run",
			},
			Builds: map[string]manifest.InnerBuilds{
				"linux": {
					Build: manifest.Build{
						Deps: []string{"gcc"},
					},
				},
				"windows": {
					Build: manifest.Build{
						Out: "main.exe",
						Env: []string{"CGO_ENABLED=1"},
						Cmd: "go build -o $OUT.exe",
					},
				},
			},
		},
	}
	assert.Equal(t, expected, *p.Build())
	assert.Equal(t, "go", p.Manifest().Runtime)
}
