package main

import (
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/install"
	"github.com/cryptopunkscc/portal/runner/observe"
	"github.com/cryptopunkscc/portal/runner/uninstall"
	apps2 "github.com/cryptopunkscc/portal/runtime/apps"
	"github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"text/tabwriter"
)

func main() { cli.Run(Application{}.cliHandler()) }

type Application struct{}

func (a Application) cliHandler() cmd.Handler {
	return cmd.Handler{
		Name: "apps",
		Desc: "Applications management.",
		Sub: cmd.Handlers{
			{
				Func: a.listApps,
				Name: "list l",
				Desc: "List installed apps.",
			},
			{
				Func: install.Runner(a.dir()).Run,
				Name: "install i",
				Desc: "Install app.",
				Params: cmd.Params{
					{Type: "string", Desc: "Path to containing directory"},
				},
			},
			{
				Func: uninstall.Runner(a.src()),
				Name: "delete d",
				Desc: "Uninstall app.",
				Params: cmd.Params{
					{Type: "string", Desc: "Application name or package name"},
				},
			},
			{
				Name: "serve s",
				Desc: "Serve apps.",
				Func: apphost.Default().Router(cmd.Handler{
					Name: "observe",
					Func: observe.NewRun(a.dir()),
				}).Run,
			},
		},
	}
}

func (a Application) dir() string        { return apps2.Dir }
func (a Application) src() target.Source { return apps2.Source }

func (a Application) listApps() (out string) {
	buffer := &bytes.Buffer{}
	w := tabwriter.NewWriter(buffer, 0, 0, 0, ' ', 0)
	t := "\t"
	for _, app := range apps.ResolveAll.List(a.src()) {
		m := app.Manifest()
		_, _ = fmt.Fprintln(w, m.Name, t, m.Title, t, m.Description, t, m.Package, t, m.Version)
	}
	w.Flush()
	return buffer.String()
}
