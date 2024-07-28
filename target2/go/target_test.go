package golang

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testManifest = &target.Manifest{
	Name:    "test",
	Title:   "test",
	Package: "test.go",
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
