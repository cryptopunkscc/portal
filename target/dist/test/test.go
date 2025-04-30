package test

import (
	_ "embed"
	"github.com/cryptopunkscc/portal/test"
	"os"
	"path/filepath"
	"testing"
)

//go:embed portal.yml
var PortalYaml []byte

func CreatePortal(t *testing.T, manifest []byte, path ...string) (dir string) {
	dir = test.CleanMkdir(t, path...)
	p := filepath.Join(dir, "portal.yml")
	err := os.WriteFile(p, manifest, 0644)
	test.AssertErr(t, err)
	return
}
