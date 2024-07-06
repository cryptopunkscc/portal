package template

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/target/template"
	"github.com/pterm/pterm"
)

func List() (err error) {
	resolve := target.Any[target.Template](target.Try(template.Resolve))
	s := source.FromFS(template.TemplatesFs)
	table := pterm.TableData{{"Short Name", "Template", "Description"}}

	for _, tt := range source.List(resolve, s) {
		t := tt.Info()
		table = append(table, []string{t.ShortName, t.Name, t.Description})
	}

	err = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return
}
