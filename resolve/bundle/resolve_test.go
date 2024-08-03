package bundle

import (
	"embed"
	"github.com/cryptopunkscc/portal/resolve/source"
	"io/fs"
	"path/filepath"
	"testing"
)

//go:embed test
var testEmbedFs embed.FS
var portalJson = "portal.json"
var testPortalPath = filepath.Join("test", "test.portal")

func TestResolve_File(t *testing.T) {
	src, err := source.File(testPortalPath)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Resolve(src)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fs.Stat(b.Files(), portalJson)
	if err != nil {
		t.Fatal(err)
	}
}

func TestResolve_Embed(t *testing.T) {
	src, err := source.Embed(testEmbedFs).Sub(testPortalPath)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Resolve(src)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fs.Stat(b.Files(), portalJson)
	if err != nil {
		t.Fatal(err)
	}
}
