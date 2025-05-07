package resolve

import (
	"fmt"
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
		target.Resolve[target.Portal_]
	}{
		{
			name: "dist",
			Resolve: target.Any[target.Portal_](
				portal.Resolve_.Try,
			),
		},
		{
			name: "js/test/project",
			Resolve: target.Any[target.Portal_](
				golang.ResolveProject.Try,
				js.ResolveDist.Try,
				js.ResolveBundle.Try,
				js.ResolveProject.Try,
				html.ResolveDist.Try,
				html.ResolveBundle.Try,
				html.ResolveProject.Try,
				exec.ResolveDist.Try,
				exec.ResolveBundle.Try,
				exec.ResolveProject.Try,
			),
		},
		{
			name: ".",
			Resolve: target.Any[target.Portal_](
				golang.ResolveProject.Try,
				js.ResolveDist.Try,
				js.ResolveBundle.Try,
				js.ResolveProject.Try,
				html.ResolveDist.Try,
				html.ResolveBundle.Try,
				html.ResolveProject.Try,
				exec.ResolveDist.Try,
				exec.ResolveBundle.Try,
				exec.ResolveProject.Try,
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := source.Dir(tt.name)
			l := tt.List(d)
			assert.NotEmpty(t, l)
			for i, p := range l {
				s := target.Sprint(p)
				s = fmt.Sprintf(" - %d. %s", i, s)
				println(s)
			}
		})
	}
}
