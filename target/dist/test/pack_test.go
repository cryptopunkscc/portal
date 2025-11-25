package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
)

func TestPack(t *testing.T) {
	for _, tt := range []struct {
		name  string
		path  []string
		clear string
	}{
		{
			name:  "default",
			clear: "build",
		},
		{
			name:  "custom path",
			path:  []string{".test_out"},
			clear: ".test_out",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.RemoveAll(tt.clear)

			s, err := source.Embed(DistFS).Sub("test_dist")
			test.AssertErr(t, err)

			d, err := dist.Resolve_(s)
			test.AssertErr(t, err)

			err = dist.Pack(d, tt.path...)
			test.AssertErr(t, err)

			n := "package_1.2.3.portal"
			p := filepath.Join("build", n)
			stat, err := os.Stat(p)
			test.AssertErr(t, err)
			assert.False(t, stat.IsDir())
			assert.Equal(t, n, stat.Name())
		})
	}
}
