package portald

import (
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
	apps2 "github.com/cryptopunkscc/portal/runtime/apps"
	"text/tabwriter"
)

type ListAppsOpts struct {
	Hidden bool `query:"hidden h" cli:"hidden h"`
}

func (s *Runner[T]) ListApps(opts ListAppsOpts) (a Apps) {
	for _, app := range apps.ResolveAll.List(apps2.Source) {
		if opts.Hidden || !app.Manifest().Hidden {
			a = append(a, app)
		}
	}
	return a
}

type Apps []target.App_

func (a Apps) MarshalCLI() string {
	b := &bytes.Buffer{}
	w := tabwriter.NewWriter(b, 4, 4, 2, ' ', 0)
	for _, app := range a {
		m := app.Manifest()
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", m.Name, m.Version, m.Title, m.Description, m.Package, m.Schema)
	}
	_ = w.Flush()
	return b.String()
}
