package test

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/source"
	sourceTest "github.com/cryptopunkscc/portal/target/source/test"
	"github.com/cryptopunkscc/portal/target/zip"
	"github.com/cryptopunkscc/portal/test"
	"testing"
)

func TestResolve(t *testing.T) {
	plog.Verbosity = 100
	zipFile := CreateTestZip(t, "test.zip")
	file2, err := source.File(zipFile)
	test.AssertErr(t, err)

	zipSource, err := zip.Resolve(file2)
	test.AssertErr(t, err)

	t.Run("sub", func(t *testing.T) {
		t.Run("foo bar", func(t *testing.T) {
			bar, err := zipSource.Sub("foo", "bar")
			test.AssertErr(t, err)
			sourceTest.AssetBar(t, bar)
		})
		t.Run("foo", func(t *testing.T) {
			foo, err := zipSource.Sub("foo")
			test.AssertErr(t, err)
			sourceTest.AssertFoo(t, foo)
			t.Run("bar", func(t *testing.T) {
				bar, err := foo.Sub("bar")
				test.AssertErr(t, err)
				sourceTest.AssetBar(t, bar)
			})
		})
	})
}
