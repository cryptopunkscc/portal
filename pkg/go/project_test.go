package golang

import (
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProject_Resolve(t *testing.T) {
	p, err := ResolveProject()
	test.AssertErr(t, err)
	assert.NotEmpty(t, p.Dir)
	assert.NotEmpty(t, p.Dir)
	assert.NotEmpty(t, p.Mod)

	d, err := p.Dependency("astrald")
	test.AssertErr(t, err)
	assert.NotEmpty(t, d.Name)
	assert.NotEmpty(t, d.Version)
	assert.NotEmpty(t, d.Replace)
}
