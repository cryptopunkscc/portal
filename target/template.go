package target

const TemplateInfoFileName = "template.json"

type TemplateInfo struct {
	Name        string `json:"name"`
	ShortName   string `json:"shortname"`
	Author      string `json:"author"`
	Description string `json:"description"`
	HelpUrl     string `json:"helpurl"`
}
