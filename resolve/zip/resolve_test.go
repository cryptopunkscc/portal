package zip

import (
	"embed"
	"github.com/cryptopunkscc/portal/resolve/source"
	"io/fs"
	"testing"
)

//go:embed embed
var embedTestData embed.FS

func TestResolve(t *testing.T) {
	var embedRoot = source.Embed(embedTestData)
	zipSrc, err := embedRoot.Sub("embed/test.zip")
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := Resolve(zipSrc)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"foo", "bar/baz"} {
		if _, err = fs.Stat(bundle.Files(), s); err != nil {
			t.Fatal("s", err)
		}
	}
}
