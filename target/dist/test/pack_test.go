package test

import (
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestPack(t *testing.T) {
	_ = os.RemoveAll("build")

	s, err := source.Embed(DistFS).Sub("test_dist")
	test.AssertErr(t, err)

	d, err := dist.Resolve_(s)
	test.AssertErr(t, err)

	err = dist.Pack(d)
	test.AssertErr(t, err)

	n := "package_version.portal"
	p := filepath.Join("build", n)
	stat, err := os.Stat(p)
	test.AssertErr(t, err)
	assert.False(t, stat.IsDir())
	assert.Equal(t, n, stat.Name())
}
