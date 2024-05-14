package template

import "github.com/pterm/pterm"

func PrintList() error {
	templates, err := Templates()
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
