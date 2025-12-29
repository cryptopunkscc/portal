package test

import (
	"context"
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
)

//go:embed apps
var AppsFS embed.FS

func TestProjectHostRunner_Run(t *testing.T) {
	config := portal.Config{}
	config.Apps = ".test_apps"

	expectedManifest := manifest.App{
		Name:        "name",
		Title:       "title",
		Description: "description",
		Package:     "package",
		Version:     1,
		Runtime:     "project_host_runner",
	}

	expectedArgs := []string{
		"run",
		expectedManifest.Runtime + "\"",
		`"apps/test app"`,
		"foo",
		"bar",
	}

	assertRunApp := func(
		ctx context.Context,
		manifest manifest.App,
		path string,
		args ...string,
	) (err error) {
		assert.Equal(t, expectedManifest, manifest)
		assert.Equal(t, "go", path)

		args[1] = filepath.Base(args[1])
		assert.Equal(t, expectedArgs, args)
		return nil
	}

	runner := exec.Runner{
		Config:     config,
		RunAppFunc: assertRunApp,
	}

	apps, err := source.Embed(AppsFS).Sub("apps")
	test.AssertErr(t, err)

	testApp, err := apps.Sub("test app")
	test.AssertErr(t, err)

	err = os.RemoveAll(config.Apps)
	test.AssertErr(t, err)

	err = apps.CopyTo(config.Apps)
	test.AssertErr(t, err)

	err = runner.ProjectHost().Run(context.Background(), testApp, "foo", "bar")
	test.AssertErr(t, err)
}
