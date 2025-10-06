package portald

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/cryptopunkscc/portal/api/target"
)

type ListAppsOpts struct {
	Hidden bool   `query:"hidden h" cli:"hidden h"`
	Scope  string `query:"scope s" cli:"scope s"`
}

func (opts ListAppsOpts) includes(app target.Portal_) bool {
	if opts.Scope == "" {
		return opts.Hidden || !app.Config().Hidden
	} else {
		appType := app.Manifest().Type
		if appType == "" {
			appType = "api"
		}
		return strings.Contains(opts.Scope, appType)
	}
}

func (s *Service) InstalledApps(opts ListAppsOpts) Apps {
	a := target.Portals[target.Portal_]{}
	for _, app := range s.Resolve.List(s.apps()) {
		if opts.includes(app) {
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
