package test

import (
	"github.com/cryptopunkscc/portal/pkg/zip"
	"github.com/cryptopunkscc/portal/target/js"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestResolveDist(t *testing.T) {
	s, err := source.Embed(distFS).Sub("dist")
	test.AssertErr(t, err)

	p, err := js.ResolveDist(s)
	test.AssertErr(t, err)

	assert.Equal(t, "js", p.Manifest().Name)
}

func TestResolveBundle(t *testing.T) {
	d := test.CleanMkdir(t, ".test_bundle")
	n := filepath.Join(d, "test_js.bundle")

	err := zip.PackFS(distFS, "dist", n)
	test.AssertErr(t, err)

	s, err := source.File(n)
	test.AssertErr(t, err)

	p, err := js.ResolveBundle(s)
	test.AssertErr(t, err)

	assert.Equal(t, "js", p.Manifest().Name)
}

func TestResolveProject(t *testing.T) {
	s, err := source.Embed(projectFS).Sub("project")
	test.AssertErr(t, err)

	p, err := js.ResolveProject(s)
	test.AssertErr(t, err)

	assert.Equal(t, "js-rollup", p.Manifest().Name)
}
