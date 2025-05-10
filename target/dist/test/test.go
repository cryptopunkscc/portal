package test

import (
	"embed"
	_ "embed"
	"github.com/cryptopunkscc/portal/pkg/test"
	"os"
	"path/filepath"
	"testing"
)

//go:embed test_dist/portal.yml
var PortalYaml []byte

//go:embed test_dist
var DistFS embed.FS

func CreatePortal(t *testing.T, manifest []byte, path ...string) (dir string) {
	dir = test.CleanMkdir(t, path...)
	p := filepath.Join(dir, "portal.yml")
	err := os.WriteFile(p, manifest, 0644)
	test.AssertErr(t, err)
	return
}
