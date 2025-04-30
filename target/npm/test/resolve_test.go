package test

import (
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveModule(t *testing.T) {
	s, err := source.Embed(moduleFS).Sub("module")
	test.AssertErr(t, err)

	p, err := npm.ResolveNodeModule(s)
	test.AssertErr(t, err)

	assert.Equal(t, "test_build", p.PkgJson().Scripts.Build)
	assert.True(t, p.PkgJson().CanBuild())
}

func TestResolveProject(t *testing.T) {
	s, err := source.Embed(projectFS).Sub("project")
	test.AssertErr(t, err)

	p, err := npm.Resolve_(s)
	test.AssertErr(t, err)

	assert.Equal(t, "js-rollup", p.Manifest().Name)
	assert.Equal(t, "test_build", p.PkgJson().Scripts.Build)
	assert.True(t, p.PkgJson().CanBuild())
}
