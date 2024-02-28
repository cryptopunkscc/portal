//go:build dev

package clir

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/create"
	"github.com/cryptopunkscc/go-astral-js/pkg/create/templates"
	"github.com/cryptopunkscc/go-astral-js/pkg/dev"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/leaanthony/clir"
	"github.com/pterm/pterm"
	"log"
)

func Run(bindings runner.Bindings) {
	cli := clir.NewCli(PortalName, PortalDevDescription, PortalVersion)
	cli.NewSubCommandFunction("create", "Create production bundle.", cliInit)
	cli.NewSubCommandFunction("dev", "Run development server for given dir.", cliDevelopment(bindings))
	cli.NewSubCommandFunction("open", "Execute app from bundle, dir, or file.", cliApplication(bindings))
	cli.NewSubCommandFunction("build", "Build application.", cliBuild)
	cli.NewSubCommandFunction("bundle", "Create production bundle.", cliBundle)
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}

func cliDevelopment(bindings runner.Bindings) func(f *FlagsPath) (err error) {
	return func(f *FlagsPath) (err error) {
		return dev.Run(f.Path, bindings)
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
	Dir      string `name:"d" description:"Project directory"`
	Force    bool   `name:"f" description:"Force recreate project"`
	List     bool   `name:"l" description:"List available templates"`
}

func cliInit(f *FlagsInit) error {
	if f.List {
		return cliList()
	} else {
		return create.Run(f.Name, f.Dir, f.Template, f.Force)
	}
}

func cliList() error {
	templateList, err := templates.List()
	if err != nil {
		return err
	}

	pterm.DefaultSection.Println("Available templates")

	table := pterm.TableData{{"Template", "Short Name", "Description"}}
	for _, template := range templateList {
		table = append(table, []string{template.Name, template.ShortName, template.Description})
	}
	err = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return err
}
