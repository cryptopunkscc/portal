package golang

import (
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

var testManifest = target.Manifest{
	Name:        "tray",
	Title:       "Portal Tray",
	Description: "Tray icon for Portal.",
	Package:     "portal.tray",
	Exec:        "main",
}

var testBuild = target.Builds{
	"default": target.Build{Exec: "main", Cmd: "go build -o dist/main"},
	"linux":   target.Build{Exec: "main", Cmd: "go build -o dist/main", Deps: []string{"gcc", "libgtk-3-dev", "libayatana-appindicator3-dev"}},
	"windows": target.Build{Exec: "main.exe", Cmd: "go build -ldflags -H=windowsgui -o dist/main.exe", Env: []string{"CGO_ENABLED=1"}},
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
	assert.Equal(t, &testManifest, bundle.Dist().Manifest())
}

func TestResolveBundle(t *testing.T) {
	name := "test/build/portal.tray_.portal"
	file, err := source.File(name)
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
	assert.Equal(t, "main", bundle.Target().Executable().Path())
	assert.True(t, strings.HasSuffix(bundle.Target().Executable().Abs(), filepath.Join(name, "main")))
}
