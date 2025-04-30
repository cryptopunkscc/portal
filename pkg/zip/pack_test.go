package zip

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test"
	"os"
	"path/filepath"
	"testing"
)

func TestPackDir(t *testing.T) {
	plog.Verbosity = 100
	srcDir := test.CleanMkdir(t, ".test_src")
	dstDir := test.CleanMkdir(t, ".test_dst")
	if err := os.WriteFile(filepath.Join(srcDir, "foo"), []byte("foo"), 0644); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(srcDir, "bar"), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "bar", "baz"), []byte("baz"), 0644); err != nil {
		panic(err)
	}

	err := Pack(srcDir, filepath.Join(dstDir, "pkg.zip"))
	if err != nil {
		plog.Println(err)
		t.FailNow()
	}
}

func TestPackFS(t *testing.T) {
	plog.Verbosity = 100
	srcDir := test.CleanMkdir(t, ".test_src")
	dstDir := test.CleanMkdir(t, ".test_dst")
	if err := os.WriteFile(filepath.Join(srcDir, "foo"), []byte("foo"), 0644); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(srcDir, "bar"), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "bar", "baz"), []byte("baz"), 0644); err != nil {
		panic(err)
	}

	err := PackFS(os.DirFS(srcDir), ".", filepath.Join(dstDir, "pkg.zip"))
	if err != nil {
		plog.Println(err)
		t.FailNow()
	}
}
