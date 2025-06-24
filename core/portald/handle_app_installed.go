package portald

import (
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"text/tabwriter"
)

type ListAppsOpts struct {
	Hidden bool `query:"hidden h" cli:"hidden h"`
}

func (s *Service) InstalledApps(opts ListAppsOpts) Apps {
	a := target.Portals[target.Portal_]{}
	for _, app := range s.Resolve.List(s.apps()) {
		if opts.Hidden || !app.Config().Hidden {
			a = append(a, app)
		}
	}
	a = a.Reduced()
	return Apps(a)
}

type Apps target.Portals[target.Portal_]

func (a Apps) MarshalCLI() string {
	b := &bytes.Buffer{}
	w := tabwriter.NewWriter(b, 4, 4, 2, ' ', 0)
	for _, app := range a {
		m := app.Manifest()
		v := ""
		if d, ok := app.(target.Dist_); ok {
			v = d.Version()
		} else {
			v = fmt.Sprintf("%d.%d", m.Version, app.Api().Version)
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", m.Name, v, m.Title, m.Description, m.Package, m.Runtime)
	}
	_ = w.Flush()
	return b.String()
}
