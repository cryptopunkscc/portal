package app

import (
	"embed"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

//go:embed test
var embedFS embed.FS

var (
	fsSrc    target.Source
	embedSrc target.Source
)

var testManifest = &target.Manifest{
	Name:    "js",
	Title:   "js",
	Package: "new.portal.js",
	Version: "0.0.0",
}

func init() {
	dir := filepath.Join("test", "app")
	var err error
	if fsSrc, err = source.File(dir); err != nil {
		panic(err)
	}
	if embedSrc, err = source.Embed(embedFS).Sub(dir); err != nil {
		panic(err)
	}
}

func TestResolve_FS(t *testing.T) {
	a, err := Resolve[any](fsSrc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, testManifest, a.Manifest())
}

func TestResolve_Embed(t *testing.T) {
	a, err := Resolve[any](embedSrc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, testManifest, a.Manifest())
}
