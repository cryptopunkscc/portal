package test

import (
	golang "github.com/cryptopunkscc/portal/target/go"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolve(t *testing.T) {
	s, err := source.Embed(goFS).Sub("go")
	test.AssertErr(t, err)

	p, err := golang.ResolveProject(s)
	test.AssertErr(t, err)

	assert.Equal(t, "go", p.Manifest().Schema)
}
