package source

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed embed
var testEmbedFs embed.FS

func TestSource_Sub_file(t *testing.T) {
	file := "embed/file"
	src, err := Embed(testEmbedFs).Sub(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, file, src.Path())
	assert.Equal(t, file, src.Abs())
	assert.False(t, src.IsDir())
}

func TestSource_Sub_dir(t *testing.T) {
	file := "embed"
	src, err := Embed(testEmbedFs).Sub(file)
	if err != nil {
		t.Error(err)
	}
	_, err = src.Files().Open("file")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, ".", src.Path())
	assert.Equal(t, file, src.Abs())
}

func TestSource_FS_file(t *testing.T) {
	dir, file := setup(t)
	defer clean(t)
	full := filepath.Join(dir, file)
	src, err := File(dir, file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, file, src.Path())
	assert.True(t,
		strings.HasSuffix(src.Abs(), full),
		"no suffix:\n%s\n%s", src.Abs(), full,
	)
}

func TestSource_FS_dir(t *testing.T) {
	dir, file := setup(t)
	defer clean(t)
	src, err := File(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = src.Files().Open(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ".", src.Path())
	assert.True(t,
		strings.HasSuffix(src.Abs(), dir),
		"no suffix:\n%s\n%s", src.Abs(), dir,
	)
}

func setup(t *testing.T) (dir string, file string) {
	dir = "test_tmp"
	file = "file"
	clean(t)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Create(filepath.Join(dir, file)); err != nil {
		t.Fatal(err)
	}
	return
}

func clean(_ *testing.T) {
	_ = os.RemoveAll("test_tmp")
}
