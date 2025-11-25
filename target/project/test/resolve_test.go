package test

import (
	"testing"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/project"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestResolve_(t *testing.T) {
	m := DevPortalYaml
	assert.NotEmpty(t, m)

	expected := manifest.App{}
	err := yaml.Unmarshal(m, &expected)
	test.AssertErr(t, err)

	sub, err := source.Embed(ProjectFS).Sub("test_project")
	test.AssertErr(t, err)

	p, err := project.Resolve_(sub)
	test.AssertErr(t, err)

	actual := *p.Manifest()

	assert.Equal(t, expected, actual)
}
