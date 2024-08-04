package golang

import (
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testManifest = target.Manifest{
	Name:    "test",
	Title:   "test",
	Package: "test.go",
	Version: "0.0.0",
}

var testBuild = target.Builds{
	"default": target.Build{Cmd: "cmd1", Env: []string(nil)},
	"windows": target.Build{Cmd: "cmd1", Env: []string{"FOO=bar"}},
}

func TestResolveProject(t *testing.T) {
	file, err := source.File("test")
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := ResolveProject(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &testManifest, bundle.Manifest())
	assert.Equal(t, testBuild, bundle.Build())
}

func TestResolveBundle(t *testing.T) {
	file, err := source.File("test/build/test.go_0.0.0.portal")
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := exec.ResolveBundle(file)
	if err != nil {
		t.Fatal(err)
	}
	bundleManifest := testManifest
	bundleManifest.Exec = "main"
	assert.Equal(t, &bundleManifest, bundle.Manifest())
	assert.Equal(t, "", bundle.Target().Executable().Path())
	assert.Equal(t, "", bundle.Target().Executable().Abs())
}
