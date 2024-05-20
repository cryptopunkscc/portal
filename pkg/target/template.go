package target

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
)

const TemplateInfoFileName = "template.json"

type TemplateInfo struct {
	Name        string `json:"name"`
	ShortName   string `json:"shortname"`
	Author      string `json:"author"`
	Description string `json:"description"`
	HelpUrl     string `json:"helpurl"`
}

func ReadTemplateInfoFS(src fs.FS) (i TemplateInfo, err error) {
	err = i.LoadFs(src, TemplateInfoFileName)
	return
}

func (m *TemplateInfo) LoadFs(src fs.FS, name string) (err error) {
	file, err := fs.ReadFile(src, name)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &m)
	return
}

func (m *TemplateInfo) LoadPath(src string, name string) (err error) {
	bytes, err := os.ReadFile(path.Join(src, name))
	if err != nil {
		return
	}
	return json.Unmarshal(bytes, &m)
}
