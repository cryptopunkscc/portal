//go:build dev

package clir

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/feat/bundle"
	"github.com/cryptopunkscc/go-astral-js/feat/create"
	"github.com/cryptopunkscc/go-astral-js/feat/dev"
	"github.com/cryptopunkscc/go-astral-js/feat/publish"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/leaanthony/clir"
	"github.com/pterm/pterm"
	"log"
	"strings"
)

func Run(ctx context.Context, bindings runtime.New) {
	cli := clir.NewCli(portal.Name, portal.DevDescription, portal.Version)

	flags := &FlagsPath{}
	cli.AddFlags(flags)
	cli.Action(func() error { return cliOpen(ctx, bindings)(flags) })

	cli.NewSubCommand("launcher", "Start portal launcher GUI.").Action(cliLauncher(ctx, bindings))
	cli.NewSubCommandFunction("create", "Create production bundle.", cliCreate)
	cli.NewSubCommandFunction("dev", "Run development server for given dir.", cliDevelopment(bindings))
	cli.NewSubCommandFunction("open", "Execute app from bundle, dir, or file.", cliOpen(ctx, bindings))
	cli.NewSubCommandFunction("build", "Build application.", cliBuild)
	cli.NewSubCommandFunction("bundle", "Create production bundle.", cliBundle)
	cli.NewSubCommandFunction("publish", "Publish bundles from given path to storage", cliPublish)
	cli.NewSubCommandFunction("install", "Install bundles from given path", cliInstall)
	cli.NewSubCommandFunction("serve", "Serve api through jrpc adapter", cliSrv(ctx, bindings))
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}

func cliDevelopment(bindings runtime.New) func(f *FlagsPath) (err error) {
	return func(f *FlagsPath) (err error) {
		return dev.Run(bindings, f.Path)
	}
}

func cliBuild(f *FlagsPath) error {
	return build.Run(f.Path)
}

func cliBundle(f *FlagsPath) error {
	return bundle.RunAll(f.Path)
}

type FlagsInit struct {
	Template string `name:"t" description:"Name of built-in template to use, path to template or template url"`
	Name     string `name:"n" description:"Name of project"`
	Force    bool   `name:"f" description:"Force recreate project"`
	List     bool   `name:"l" description:"List available templates"`
	Dir      string `pos:"1" description:"Project directory"`
}

var emptyFlagsInit = FlagsInit{}

func cliCreate(f *FlagsInit) error {
	switch {
	case *f == emptyFlagsInit:
		return cliList()
	case f.List:
		return cliList()
	default:
		temps := strings.Split(f.Template, " ")
		return create.Run(f.Name, f.Dir, temps, f.Force)
	}
}

func cliList() error {
	templates, err := project.Templates()
	if err != nil {
		return err
	}

	pterm.DefaultSection.Println("Available templates")

	table := pterm.TableData{{"Template", "Short Name", "Description"}}
	for _, template := range templates {
		table = append(table, []string{template.Name, template.ShortName, template.Description})
	}
	err = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return err
}

func cliPublish(f *FlagsPath) error {
	return publish.Run(f.Path)
}

func cliInstall(f *FlagsPath) error {
	return appstore.Install(f.Path)
}

func cliSrv(ctx context.Context, bindings runtime.New) func(*struct{}) error {
	return func(_ *struct{}) error {
		return serve.Run(ctx, bindings)
	}
}
