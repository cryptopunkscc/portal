package source

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSource_Sub_file(t *testing.T) {
	file := "test/file"
	src, err := Embed(testEmbedFs).Sub(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, file, src.Path())
	assert.Equal(t, "", src.Abs())
}

func TestSource_Sub_dir(t *testing.T) {
	file := "test"
	src, err := Embed(testEmbedFs).Sub(file)
	if err != nil {
		t.Error(err)
	}
	_, err = src.Files().Open("file")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, ".", src.Path())
	assert.Equal(t, "", src.Abs())
}

func TestSource_FS_file(t *testing.T) {
	file := "test/file"
	src, err := File(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "file", src.Path())
	assert.True(t, strings.HasSuffix(src.Abs(), file), "no suffix:\n%s\n%s", src.Abs(), file)
}

func TestSource_FS_dir(t *testing.T) {
	file := "test"
	src, err := File(file)
	if err != nil {
		t.Fatal(err)
	}
	_, err = src.Files().Open("file")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ".", src.Path())
	assert.True(t, strings.HasSuffix(src.Abs(), file), "no suffix:\n%s\n%s", src.Abs(), file)
}

//go:embed test
var testEmbedFs embed.FS
