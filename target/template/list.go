package template

import (
	"bytes"
	"embed"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/pterm/pterm"
)

func ListFrom(fs embed.FS) Templates {
	src := source.Embed(fs)
	return Resolve.List(src)
}

type Templates []target.Template

func (t Templates) MarshalCLI() string {
	if len(t) == 0 {
		return "no templates found"
	}

	table := pterm.TableData{{"Short Name", "Template", "Description"}}
	buffer := bytes.NewBufferString("")

	for _, tt := range t {
		info := tt.Info()
		row := []string{info.ShortName, info.Name, info.Description}
		table = append(table, row)
	}

	err := pterm.DefaultTable.
		WithHasHeader(true).
		//WithBoxed(true).
		WithData(table).
		WithWriter(buffer).
		Render()
	if err != nil {
		return err.Error()
	} else {
		return buffer.String()
	}
}
