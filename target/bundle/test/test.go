package test

import (
	"path/filepath"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/pkg/zip"
	disttest "github.com/cryptopunkscc/portal/target/dist/test"
)

func CreateBundleM(t *testing.T, manifest []byte, dst ...string) string {
	src := disttest.CreatePortal(t, manifest)
	return CreateBundle(t, src, dst...)
}

func CreateBundle(t *testing.T, src string, dst ...string) string {
	dir, file := filepath.Split(filepath.Join(dst...))
	testDst := test.CleanMkdir(t, dir)
	testZip := filepath.Join(testDst, file)
	if err := zip.Pack(src, testZip); err != nil {
		panic(err)
	}
	return testZip
}
