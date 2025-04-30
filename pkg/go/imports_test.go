package golang

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"path/filepath"
	"testing"
)

func Test_ListImports(t *testing.T) {
	src, _ := FindProjectRoot()
	src = filepath.Join(src, "pkg/go/imports.go")
	imports, err := Imports(src)
	if err != nil {
		plog.Println(err)
	}

	for i, s := range imports {
		t.Log(i, s)
	}
}
