package test

import (
	"testing"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/bundle"
	disttest "github.com/cryptopunkscc/portal/target/dist/test"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestResolve(t *testing.T) {
	m := disttest.PortalYaml
	assert.NotEmpty(t, m)

	expected := manifest.App{}
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
