package resolve

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/exec"
	golang "github.com/cryptopunkscc/portal/target/go"
	"github.com/cryptopunkscc/portal/target/html"
	"github.com/cryptopunkscc/portal/target/js"
	"github.com/cryptopunkscc/portal/target/portal"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	plog.Verbosity = 100
}

func TestResolve_List(t *testing.T) {
	for _, tt := range []struct {
		name string
		target.Resolve[target.Source]
	}{
		{
			name:    "dist",
			Resolve: target.Try(portal.Resolve_),
		},
		{
			name: "js/test/project",
			Resolve: target.Any[target.Source](
				target.Try(golang.ResolveProject),
				target.Try(js.ResolveDist),
				target.Try(js.ResolveBundle),
				target.Try(js.ResolveProject),
				target.Try(html.ResolveDist),
				target.Try(html.ResolveBundle),
				target.Try(html.ResolveProject),
				target.Try(exec.ResolveDist),
				target.Try(exec.ResolveBundle),
				target.Try(exec.ResolveProject),
			),
		},
		{
			name: ".",
			Resolve: target.Any[target.Source](
				target.Try(golang.ResolveProject),
				target.Try(js.ResolveDist),
				target.Try(js.ResolveBundle),
				target.Try(js.ResolveProject),
				target.Try(html.ResolveDist),
				target.Try(html.ResolveBundle),
				target.Try(html.ResolveProject),
				target.Try(exec.ResolveDist),
				target.Try(exec.ResolveBundle),
				target.Try(exec.ResolveProject),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := source.Dir(tt.name)
			l := tt.List(d)
			assert.NotEmpty(t, l)
			for i, p := range l {
				plog.Printf("%d %T %v", i, p, p)
			}
		})
	}
}
