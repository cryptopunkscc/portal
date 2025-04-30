package test

import (
	"bytes"
	_ "embed"
	"github.com/cryptopunkscc/portal/api/target"
	test2 "github.com/cryptopunkscc/portal/resolve/bundle/test"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestResolveManifest(t *testing.T) {
	open, err := os.Open("portal.yml")
	test.AssertErr(t, err)
	m := target.Manifest{}

	_, err = m.ReadFrom(open)
	test.AssertErr(t, err)
}

func TestResolveDist(t *testing.T) {
	p := CreateDistExec(t, ".test_dist")

	expected := target.Manifest{}
	b := bytes.NewBuffer(PortalYaml)
	_, err := expected.ReadFrom(b)
	test.AssertErr(t, err)

	s := source.Dir(p)
	d, err := exec.ResolveDist(s)
	test.AssertErr(t, err)

	actual := *d.Manifest()
	assert.Equal(t, expected, actual)
}

func TestResolveBundle(t *testing.T) {
	p := CreateDistExec(t, ".test_dist")
	p = test2.CreateBundle(t, p, ".test_bundle", "test_portal.bundle")

	expected := target.Manifest{}
	b := bytes.NewBuffer(PortalYaml)
	_, err := expected.ReadFrom(b)
	test.AssertErr(t, err)

	s, err := source.File(p)
	test.AssertErr(t, err)

	d, err := exec.ResolveBundle(s)
	test.AssertErr(t, err)

	actual := *d.Manifest()
	assert.Equal(t, expected, actual)
}

func TestResolveProject(t *testing.T) {
	d := CreateProjectExec(t, ".test_project")
	s := source.Dir(d)

	p, err := exec.ResolveProject(s)
	test.AssertErr(t, err)

	assert.Equal(t, "exec", p.Manifest().Exec)
	assert.NotZero(t, d)
}
