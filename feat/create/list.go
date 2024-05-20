package create

import (
	. "github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/source"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/template"
	"github.com/pterm/pterm"
)

func List() (err error) {
	resolve := Any[Template](Try(template.Resolve))
	s := source.Resolve(template.TemplatesFs)
	table := pterm.TableData{{"Short Name", "Template", "Description"}}

	for tt := range source.Stream(resolve, s) {
		t := tt.Info()
		table = append(table, []string{t.ShortName, t.Name, t.Description})
	}

	err = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return
}
