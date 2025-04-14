package target

import (
	"github.com/cryptopunkscc/portal/pkg/dec"
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"strings"
	"testing"
)

func TestTest(t *testing.T) {
	assert.Equal(t, []string{"", "a"}, strings.Split("-a", "-"))
}

func Test_mapped_loadBuilds(t *testing.T) {
	allUnmarshalers := all.Unmarshalers
	t.Cleanup(func() { all.Unmarshalers = allUnmarshalers })
	all.Unmarshalers = dec.Load(func(dst any, src fs.FS, name string) (err error) {
		dst.(*builds).Builds = Builds{
			"default": Build{
				Out:  "out1",
				Deps: []string{"dep1"},
				Env:  []string{"env1"},
			},
			"os": Build{
				Out:  "out2",
				Deps: []string{"dep2"},
				Env:  []string{"env2"},
			},
			"os-arch1-arch2": Build{
				Out:  "out3",
				Deps: []string{"dep3"},
				Env:  []string{"env3"},
			},
		}
		return
	})
	m := mapped{}
	m.loadBuilds(nil)
	f := m.flatten()
	t.Log(f)
}
