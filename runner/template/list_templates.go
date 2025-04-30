package template

import (
	"bytes"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/target/template"
	"github.com/pterm/pterm"
)

func Map() (templates map[string]target.Template) {
	templates = map[string]target.Template{}
	for _, t := range List() {
		templates[t.Name()] = t
	}
	return
}

func List() Templates {
	src := source.Embed(template.TemplatesFs)
	return template.Resolve.List(src)
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
