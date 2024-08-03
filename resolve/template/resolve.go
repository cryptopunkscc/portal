package template

import (
	"encoding/json"
	"errors"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
	"path/filepath"
)

type source struct {
	target.Source
	info target.TemplateInfo
}

func (t *source) Info() target.TemplateInfo { return t.info }
func (t *source) Name() (name string) {
	name = t.info.ShortName
	if name == "" {
		name = filepath.Base(t.Abs())
	}
	return
}

func Resolve(src target.Source) (template target.Template, err error) {
	if !src.IsDir() {
		return nil, ErrNotTemplate
	}
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

var ErrNotTemplate = errors.New("not a template")
