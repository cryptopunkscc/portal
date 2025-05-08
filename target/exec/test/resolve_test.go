package test

import (
	_ "embed"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/target"
	test2 "github.com/cryptopunkscc/portal/target/bundle/test"
	"github.com/cryptopunkscc/portal/target/exec"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestResolveManifest(t *testing.T) {
	open, err := os.Open("portal.yml")
	test.AssertErr(t, err)
	m := manifest.App{}

	_, err = m.ReadFrom(open)
	test.AssertErr(t, err)
}

func TestResolveProject(t *testing.T) {
	expected := manifest.Dev{}
	err := expected.UnmarshalFrom(DevPortalYaml)
	test.AssertErr(t, err)

	s, err := source.Embed(ProjectFS).Sub("test_project")
	test.AssertErr(t, err)

	p, err := exec.ResolveProject(s)
	test.AssertErr(t, err)

	assert.Equal(t, expected.App, *p.Manifest())
	assert.Equal(t, expected.Api, *p.Api())
	assert.Equal(t, expected.Config, *p.Config())

	actualTarget := p.Build().Get("linux").Target
	assert.Equal(t, "out", actualTarget.Exec)
	assert.Equal(t, "linux", actualTarget.OS)
	assert.NotZero(t, actualTarget.Arch)
}

func TestResolveDist(t *testing.T) {
	expected := ExpectedDist(t)

	p := CreateDistExec(t, ".test_dist")
	s := source.Dir(p)
	actual, err := exec.ResolveDist(s)
	test.AssertErr(t, err)

	AssertDist(t, expected, actual)
}

func TestResolveBundle(t *testing.T) {
	expected := ExpectedDist(t)

	p := CreateDistExec(t, ".test_dist")
	p = test2.CreateBundle(t, p, ".test_bundle", "test_portal.bundle")
	s, err := source.File(p)
	test.AssertErr(t, err)

	actual, err := exec.ResolveBundle(s)
	test.AssertErr(t, err)

	AssertDist(t, expected, actual)
}

func ExpectedDist(t *testing.T) (expected manifest.Dist) {
	err := expected.UnmarshalFrom(PortalYaml)
	test.AssertErr(t, err)
	expected.Target.Exec = "exec"
	return
}

func AssertDist(t *testing.T, expected manifest.Dist, actual target.DistExec) {
	assert.Equal(t, expected.App, *actual.Manifest())
	assert.Equal(t, expected.Api, *actual.Api())
	assert.Equal(t, expected.Config, *actual.Config())
	assert.Equal(t, expected.Release, *actual.Release())
	assert.Equal(t, expected.Target, actual.Runtime().Target())
}
