package js

import (
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

var testManifest = &target.Manifest{
	Name:    "project",
	Title:   "project",
	Package: "new.portal.project",
	Version: "0.0.0",
}

func TestResolveBundle(t *testing.T) {
	file, err := source.File("test", "build", "new.portal.project_0.0.0.portal")
	if err != nil {
		t.Error(err)
	}
	bundle, err := ResolveBundle(file)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, testManifest, bundle.Manifest())
}

func TestResolveDist(t *testing.T) {
	file, err := source.File("test", "dist")
	if err != nil {
		t.Fatal(err)
	}
	dist, err := ResolveDist(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testManifest, dist.Manifest())
	assert.True(t, strings.HasSuffix(dist.Abs(), filepath.Join("js", "test", "dist")))
}

var testBuild = target.Builds{
	"default": target.Build{Cmd: "cmd2", Deps: []string{"dep2"}},
	"linux":   target.Build{Cmd: "cmd3", Deps: []string{"dep2", "dep3"}},
	"windows": target.Build{Cmd: "cmd4", Deps: []string{"dep2", "dep4"}},
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
	assert.Equal(t, testManifest, bundle.Manifest())
	assert.Equal(t, testBuild, bundle.Build())
}
