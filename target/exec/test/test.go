package test

import (
	_ "embed"
	disttest "github.com/cryptopunkscc/portal/target/dist/test"
	projectTest "github.com/cryptopunkscc/portal/target/project/test"
	"github.com/cryptopunkscc/portal/test"
	"os"
	"path/filepath"
	"testing"
)

//go:embed portal.yml
var PortalYaml []byte

//go:embed dev.portal.yml
var DevPortalYaml []byte

func CreateDistExec(t *testing.T, path ...string) (dir string) {
	dir = disttest.CreatePortal(t, PortalYaml, path...)
	main := filepath.Join(dir, "exec")
	c := []byte("#!/bin/bash")
	err := os.WriteFile(main, c, 0777)
	test.AssertErr(t, err)
	return
}

func CreateProjectExec(t *testing.T, path ...string) (dir string) {
	dir = projectTest.CreateProject(t, DevPortalYaml, path...)
	return
}
