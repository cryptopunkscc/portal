package golang

import (
	"log"
	"os"
	"path"
	"testing"
)

func TestWatchCache_AddFile(t *testing.T) {
	src, _ := os.Getwd()
	wd, _ := FindProjectRoot(src)
	target := path.Join(wd, "pkg/go/imports.go")
	w := NewWatchCache(wd, "github.com/cryptopunkscc/portal/")
	w.AddFile(target)
	for s, i := range w.dirs {
		log.Println(s, i)
	}
	for s, strings := range w.files {
		log.Println(s, strings)
	}
}
