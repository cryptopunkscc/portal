//go:build dev

package clir

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/feat/bundle"
	"github.com/cryptopunkscc/go-astral-js/feat/create"
	"github.com/cryptopunkscc/go-astral-js/feat/dev"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/publish"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/feat/tray"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/template"
	"github.com/leaanthony/clir"
	"github.com/pterm/pterm"
	"log"
	"strings"
)

func Run(ctx context.Context, bindings runtime.New) {
	cli := Cli{
		Cli:      clir.NewCli(portal.Name, portal.DevDescription, portal.Version),
		ctx:      ctx,
		bindings: bindings,
	}

	flags := &FlagsOpen{}
	cli.AddFlags(flags)
	cli.Action(func() error { return cli.Open(flags) })

	cli.NewSubCommand("launcher", "Start portal launcher GUI.").Action(cli.Launcher)
	cli.NewSubCommandFunction("create", "Create production bundle.", cli.Create)
	cli.NewSubCommandFunction("dev", "Run development server for given dir.", cli.Development)
	cli.NewSubCommandFunction("open", "Execute app from bundle, dir, or file.", cli.Open)
	cli.NewSubCommandFunction("build", "Build application.", cli.Build)
	cli.NewSubCommandFunction("bundle", "Create production bundle.", cli.Bundle)
	cli.NewSubCommandFunction("publish", "Publish bundles from given path to storage", cli.Publish)
	cli.NewSubCommandFunction("install", "Install bundles from given path", cli.Install)
	cli.NewSubCommandFunction("serve", "Serve api through rpc adapter", cli.Srv)
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}

type FlagsInit struct {
	Template string `name:"t" description:"Name of built-in template to use, path to template or template url"`
	Name     string `name:"n" description:"Name of project"`
	Force    bool   `name:"f" description:"Force recreate project"`
	List     bool   `name:"l" description:"List available templates"`
	Dir      string `pos:"1" description:"Project directory"`
}

var emptyFlagsInit = FlagsInit{}

type FlagsSrv struct {
	Tray bool `name:"t" description:"Display system tray indicator."`
}

func (c Cli) Development(f *FlagsOpen) (err error) {
	return dev.Run(c.ctx, c.bindings, f.Path, f.Attach)
}

func (c Cli) Build(f *FlagsPath) error {
	return build.Run(f.Path)
}

func (c Cli) Bundle(f *FlagsPath) error {
	return bundle.RunAll(f.Path)
}

func (c Cli) Create(f *FlagsInit) error {
	switch {
	case *f == emptyFlagsInit:
		return c.List()
	case f.List:
		return c.List()
	default:
		temps := strings.Split(f.Template, " ")
		return create.Run(f.Name, f.Dir, temps, f.Force)
	}
}

func (c Cli) List() error {
	templates, err := template.Templates()
	if err != nil {
		return err
	}

	pterm.DefaultSection.Println("Available templates")

	table := pterm.TableData{{"Template", "Short Name", "Description"}}
	for _, t := range templates {
		table = append(table, []string{t.Name, t.ShortName, t.Description})
	}
	err = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return err
}

func (c Cli) Publish(f *FlagsPath) error {
	return publish.Run(f.Path)
}

func (c Cli) Install(f *FlagsPath) error {
	return appstore.Install(f.Path)
}

func (c Cli) Srv(f *FlagsSrv) error {
	var t runtime.Tray
	if f.Tray {
		t = tray.Run
	}
	return serve.Run(c.ctx, "portal", open.Handlers, t)
}
