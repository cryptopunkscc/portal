package test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
)

func CreateTestDir(path ...string) string {
	testDir := filepath.Join(path...)
	if len(testDir) == 0 {
		testDir = ".test"
	}
	if err := os.RemoveAll(testDir); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(testDir, "foo"), 0777); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "foo", "bar"), []byte("baz"), 0777); err != nil {
		panic(err)
	}
	return testDir
}

func AssertFoo(t *testing.T, source target.Source) {
	assert.True(t, source.IsDir())
	file, err := source.File()
	test.AssertErr(t, err)
	stat, err := file.Stat()
	test.AssertErr(t, err)
	assert.Equal(t, stat.Name(), "foo")
	assert.True(t, stat.IsDir())
}

func AssetBar(t *testing.T, source target.Source) {
	assert.False(t, source.IsDir())
	file, err := fs.ReadFile(source.FS(), source.Path())
	test.AssertErr(t, err)
	assert.Equal(t, "baz", string(file))
}
