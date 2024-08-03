package template

import (
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/template"
	"github.com/cryptopunkscc/portal/target"
	"github.com/pterm/pterm"
)

func List() (err error) {
	table := pterm.TableData{{"Short Name", "Template", "Description"}}
	for _, tt := range target.List(
		template.Resolve,
		source.Embed(template.TemplatesFs),
	) {
		t := tt.Info()
		table = append(table, []string{t.ShortName, t.Name, t.Description})
	}
	err = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return
}
