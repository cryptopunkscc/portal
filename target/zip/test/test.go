package test

import (
	"github.com/cryptopunkscc/portal/pkg/zip"
	sourceTest "github.com/cryptopunkscc/portal/target/source/test"
	"github.com/cryptopunkscc/portal/test"
	"path/filepath"
	"testing"
)

func CreateTestZip(t *testing.T, name string, src ...string) string {
	testSrc := sourceTest.CreateTestDir(src...)
	testDst := test.CleanMkdir(t, ".test_dst")
	testZip := filepath.Join(testDst, name)
	if err := zip.Pack(testSrc, testZip); err != nil {
		panic(err)
	}
	return testZip
}

func CreateZip(t *testing.T, src, name string) string {
	testDst := test.CleanMkdir(t, ".test_dst")
	testZip := filepath.Join(testDst, name)
	if err := zip.Pack(src, testZip); err != nil {
		panic(err)
	}
	return testZip
}
