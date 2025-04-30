package test

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	disttest "github.com/cryptopunkscc/portal/resolve/dist/test"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestResolve(t *testing.T) {
	m := disttest.PortalYaml
	assert.NotEmpty(t, m)

	expected := target.Manifest{}
	err := yaml.Unmarshal(m, &expected)
	test.AssertErr(t, err)

	b := CreateBundleM(t, m, ".test_dst", "test_portal.bundle")
	f, err := source.File(b)
	test.AssertErr(t, err)

	s, err := bundle.Resolve_(f)
	test.AssertErr(t, err)

	actual := *s.Manifest()
	assert.Equal(t, expected, actual)
}
