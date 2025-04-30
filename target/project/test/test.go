package test

import (
	_ "embed"
	"github.com/cryptopunkscc/portal/test"
	"os"
	"path/filepath"
	"testing"
)

//go:embed dev.portal.yml
var DevPortalYaml []byte

func CreateProject(t *testing.T, manifest []byte, path ...string) (dir string) {
	dir = test.CleanMkdir(t, path...)
	p := filepath.Join(dir, "dev.portal.yml")
	err := os.WriteFile(p, manifest, 0644)
	test.AssertErr(t, err)
	return
}
