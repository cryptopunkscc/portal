package test

import (
	"embed"
	_ "embed"
	"github.com/cryptopunkscc/portal/api/manifest"
	"github.com/cryptopunkscc/portal/pkg/test"
	disttest "github.com/cryptopunkscc/portal/target/dist/test"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

//go:embed test_project
var ProjectFS embed.FS

//go:embed test_project/dev.portal.yml
var DevPortalYaml []byte

var PortalYaml []byte

func init() {
	dist := manifest.Dist{
		App: manifest.App{
			Name:        "name",
			Title:       "title",
			Description: "description",
			Package:     "package",
			Version:     1,
		},
	}
	switch runtime.GOOS {
	case "windows":
		dist.Target.Exec = "exec.cmd"
	default:
		dist.Target.Exec = "exec"
	}
	var err error
	PortalYaml, err = yaml.Marshal(dist)
	if err != nil {
		panic(err)
	}
}

func CreateDistExec(t *testing.T, path ...string) (dir string) {
	dir = disttest.CreatePortal(t, PortalYaml, path...)
	createExecFile(t, dir, "exec", "#!/bin/bash")
	createExecFile(t, dir, "exec.cmd", `echo hello!!!`)
	return
}

func createExecFile(t *testing.T, dir, name string, payload string) {
	main := filepath.Join(dir, name)
	err := os.WriteFile(main, []byte(payload), 0777)
	test.AssertErr(t, err)
}
