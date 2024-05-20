package project

import (
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"path"
)

type Template struct {
	target.Source
	info target.TemplateInfo
}

var _ target.Template = (*Template)(nil)

func (t *Template) Name() (name string) {
	name = t.info.ShortName
	if name == "" {
		name = path.Base(t.Abs())
	}
	return
}

var _ target.Template = (*Template)(nil)

var ErrNotTemplate = errors.New("not a template")

func (t *Template) Info() target.TemplateInfo {
	return t.info
}

func ResolveTemplate(m target.Source) (t target.Template, err error) {
	if m.IsFile() {
		return nil, ErrNotTemplate
	}
	m = m.Lift()
	info, err := target.ReadTemplateInfoFS(m.Files())
	if err != nil {
		return
	}
	t = &Template{
		Source: m,
		info:   info,
	}
	return
}
