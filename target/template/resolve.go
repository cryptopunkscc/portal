package template

import (
	"encoding/json"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"io/fs"
)

var Resolve target.Resolve[target.Template] = resolve

func resolve(src target.Source) (template target.Template, err error) {
	if !src.IsDir() {
		return nil, ErrNotTemplate
	}
	info, err := readTemplateInfo(src.FS())
	if err != nil {
		return
	}
	template = &Source{
		Source: src,
		info:   info,
	}
	return
}

func readTemplateInfo(src fs.FS) (info target.TemplateInfo, err error) {
	file, err := fs.ReadFile(src, target.TemplateInfoFileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &info)
	return
}

var ErrNotTemplate = errors.New("not a template")
