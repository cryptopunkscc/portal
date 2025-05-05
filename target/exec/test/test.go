package test

import (
	"embed"
	_ "embed"
	disttest "github.com/cryptopunkscc/portal/target/dist/test"
	"github.com/cryptopunkscc/portal/test"
	"os"
	"path/filepath"
	"testing"
)

//go:embed test_project
var ProjectFS embed.FS

//go:embed test_project/dev.portal.yml
var DevPortalYaml []byte

//go:embed portal.yml
var PortalYaml []byte

func CreateDistExec(t *testing.T, path ...string) (dir string) {
	dir = disttest.CreatePortal(t, PortalYaml, path...)
	main := filepath.Join(dir, "exec")
	c := []byte("#!/bin/bash")
	err := os.WriteFile(main, c, 0777)
	test.AssertErr(t, err)
	return
}
