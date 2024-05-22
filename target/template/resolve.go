package template

import (
	"encoding/json"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
	"io/fs"
)

var ErrNotTemplate = errors.New("not a template")

func Resolve(m target.Source) (t target.Template, err error) {
	if m.IsFile() {
		return nil, ErrNotTemplate
	}
	m = m.Lift()
	info, err := readTemplateInfo(m.Files())
	if err != nil {
		return
	}
	t = &source{
		Source: m,
		info:   info,
	}
	return
}

func readTemplateInfo(src fs.FS) (i target.TemplateInfo, err error) {
	file, err := fs.ReadFile(src, target.TemplateInfoFileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &i)
	return
}
