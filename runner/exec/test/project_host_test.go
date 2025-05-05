package test

import (
	"context"
	"embed"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

//go:embed app
var AppFS embed.FS

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
		expectedManifest.Runtime,
		"app",
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

	ctx := context.Background()
	s, err := source.Embed(AppFS).Sub("app")
	test.AssertErr(t, err)

	ss, err := source.Embed(AppsFS).Sub("apps")
	test.AssertErr(t, err)
	test.Copy(ss, config.Apps)

	err = runner.ProjectHost().Run(ctx, s, "foo", "bar")
	test.AssertErr(t, err)
}
