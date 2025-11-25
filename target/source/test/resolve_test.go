package test

import (
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/source"
)

func Test(t *testing.T) {
	testDir := CreateTestDir()

	t.Run("file", func(t *testing.T) {
		t.Run("file", func(t *testing.T) {
			sub, err := source.File(testDir, "foo", "bar")
			test.AssertErr(t, err)
			AssetBar(t, sub)
		})

		t.Run("dir", func(t *testing.T) {
			sub, err := source.File(testDir, "foo")
			test.AssertErr(t, err)
			AssertFoo(t, sub)
		})
	})

	t.Run("dir", func(t *testing.T) {
		dir := source.Dir(testDir)

		t.Run("sub", func(t *testing.T) {
			t.Run("foo bar baz", func(t *testing.T) {
				sub, err := dir.Sub("foo", "bar")
				test.AssertErr(t, err)
				AssetBar(t, sub)
			})

			t.Run("foo", func(t *testing.T) {
				sub, err := dir.Sub("foo")
				test.AssertErr(t, err)
				AssertFoo(t, sub)

				t.Run("bar", func(t *testing.T) {
					sub, err := sub.Sub("bar")
					test.AssertErr(t, err)
					AssetBar(t, sub)
				})
			})
		})
	})
}
