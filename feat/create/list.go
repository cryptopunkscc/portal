package create

import (
	. "github.com/cryptopunkscc/go-astral-js/pkg/target"
	template2 "github.com/cryptopunkscc/go-astral-js/pkg/target/template"
	"github.com/cryptopunkscc/go-astral-js/pkg/template"
	"github.com/pterm/pterm"
)

func List() (err error) {
	resolve := Any[Template](Try(template2.Resolve))
	source := NewModuleFS(template.TemplatesFs)
	table := pterm.TableData{{"Short Name", "Template", "Description"}}

	for tt := range Stream(resolve, source) {
		t := tt.Info()
		table = append(table, []string{t.ShortName, t.Name, t.Description})
	}

	err = pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return
}
