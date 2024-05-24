package template

import (
	"encoding/json"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
	"io/fs"
)

var ErrNotTemplate = errors.New("not a template")

func Resolve(src target.Source) (template target.Template, err error) {
	if src.IsFile() {
		return nil, ErrNotTemplate
	}
	src = src.Lift()
	info, err := readTemplateInfo(src.Files())
	if err != nil {
		return
	}
	template = &source{
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
