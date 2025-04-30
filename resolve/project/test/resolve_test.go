package test

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/project"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestResolve_(t *testing.T) {
	m := DevPortalYaml
	assert.NotEmpty(t, m)

	expected := target.Manifest{}
	err := yaml.Unmarshal(m, &expected)
	test.AssertErr(t, err)

	dir := CreateProject(t, m)
	s := source.Dir(dir)

	p, err := project.Resolve_(s)
	test.AssertErr(t, err)

	actual := *p.Manifest()
	assert.Equal(t, "exec", actual.Exec)

	actual.Exec = ""
	assert.Equal(t, expected, actual)
}
