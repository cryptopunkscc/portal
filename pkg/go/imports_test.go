package golang

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
	"path"
	"testing"
)

func Test_ListImports(t *testing.T) {
	plog.ErrorStackTrace = true
	src, _ := os.Getwd()
	src, _ = FindProjectRoot(src)
	src = path.Join(src, "pkg/go/imports.go")
	imports, err := Imports(src)
	if err != nil {
		plog.Println(err)
	}

	for i, s := range imports {
		t.Log(i, s)
	}
}
