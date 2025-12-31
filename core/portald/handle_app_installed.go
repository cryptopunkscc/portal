package portald

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
)

type ListAppsOpts struct {
	Hidden bool   `query:"hidden h" cli:"hidden h"`
	Scope  string `query:"scope s" cli:"scope s"`
}

func (opts ListAppsOpts) includes(app app.Metadata) bool {
	if opts.Scope == "" {
		return opts.Hidden || !app.Config.Hidden
	} else {
		appType := app.Manifest.Type
		if appType == "" {
			appType = "api"
		}
		return strings.Contains(opts.Scope, appType)
	}
}

func (s *Service) InstalledApps(opts ListAppsOpts) (a Apps) {
	m := sig.Map[string, app.Dist]{}
	for _, app := range source.CollectIt(s.appsRef(), &app.Bundle{}) {
		if opts.includes(app.Metadata) {
			m.Set(app.Metadata.Package, app.Dist)
		}
	}
	return m.Values()
}

type Apps []app.Dist

func (a Apps) MarshalCLI() string {
	b := &bytes.Buffer{}
	w := tabwriter.NewWriter(b, 4, 4, 2, ' ', 0)
	for _, app := range a {
		m := app.GetDist().Metadata
		v := fmt.Sprintf("%d.%d", m.Version, m.Api.Version)
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", m.Name, v, m.Title, m.Description, m.Package, m.Runtime)
	}
	_ = w.Flush()
	return b.String()
}
