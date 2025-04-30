package test

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestResolve(t *testing.T) {
	m := PortalYaml
	assert.NotEmpty(t, m)

	expected := target.Manifest{}
	err := yaml.Unmarshal(PortalYaml, &expected)
	test.AssertErr(t, err)

	dir := CreatePortal(t, m)
	s := source.Dir(dir)

	p, err := dist.Resolve_(s)
	test.AssertErr(t, err)
	assert.Equal(t, expected, *p.Manifest())
}
