package apps

import (
	"context"
	apps2 "github.com/cryptopunkscc/portal/api/apps"
	"github.com/cryptopunkscc/portal/test"
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const testAppsDir = "test_apps_dir"

func TestDefault_Base(t *testing.T) {
	defer test.Clean()
	defer os.RemoveAll(testAppsDir)
	src := test.Copy(test.EmbedSvelteBundle)
	apps := Dir(testAppsDir)
	ctx := context.Background()

	// observe
	var app apps2.App
	observe, err := apps.Observe(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// install
	if err = apps.Install(ctx, src.File()); err != nil {
		t.Fatal(err)
	}

	// observe create
	app = <-observe
	assert.Equal(t, fsnotify.Create, app.Event.Op)
	assert.Equal(t, test.EmbedSvelteManifest, app.Manifest())

	// list
	list, err := apps.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, list)

	// get
	got, err := apps.Get(ctx, test.EmbedSvelteManifest.Package)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, test.EmbedSvelteManifest, got.Manifest())

	// uninstall
	err = apps.Uninstall(ctx, test.EmbedSvelteManifest.Package)
	if err != nil {
		t.Fatal(err)
	}

	// observe remove
	app = <-observe
	assert.Equal(t, fsnotify.Remove, app.Event.Op)
	assert.Equal(t, test.EmbedSvelteManifest, app.Manifest())
}

func TestDefault_InstallFromPath(t *testing.T) {
	defer test.Clean()
	defer os.RemoveAll(testAppsDir)
	apps := Dir(testAppsDir)
	ctx := context.Background()

	// install
	if err := apps.InstallFromPath(ctx, "../../example"); err != nil {
		t.Fatal(err)
	}
}
