package html

import (
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testManifest = &target.Manifest{
	Name:    "test",
	Title:   "test",
	Package: "new.portal.test",
	Version: "0.0.0",
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
}

func TestResolveBundle(t *testing.T) {
	file, err := source.File("test", "build", "new.portal.test_0.0.0.portal")
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
}
