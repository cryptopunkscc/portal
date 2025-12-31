package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func init() {
	plog.Verbosity = 100
}

var DefaultTestDir = ".test"

func CleanDir(t *testing.T, path ...string) string {
	Clean(path...)
	return Dir(t, path...)
}

func Dir(t *testing.T, path ...string) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := DefaultTestDir
	if len(path) > 0 {
		dir = filepath.Join(path...)
	}
	return filepath.Join(wd, dir)
}

func Mkdir(t *testing.T, path ...string) (d string) {
	d = Dir(t, path...)
	if err := os.MkdirAll(d, 0755); err != nil {
		t.Fatal(err)
	}
	return
}

func CleanMkdir(t *testing.T, path ...string) (d string) {
	Clean(path...)
	return Mkdir(t, path...)
}

func Clean(path ...string) {
	dir := DefaultTestDir
	if len(path) > 0 {
		dir = filepath.Join(path...)
	}
	_ = os.RemoveAll(dir)
}

func NoError(t *testing.T, err error) {
	AssertErr(t, err)
}

func AssertErr(t *testing.T, err error) {
	if err != nil {
		plog.Println(err)
		t.FailNow()
	}
}
