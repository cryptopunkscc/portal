package template

import (
	"github.com/cryptopunkscc/portal/target"
	"path"
)

type source struct {
	target.Source
	info target.TemplateInfo
}

func (t *source) Name() (name string) {
	name = t.info.ShortName
	if name == "" {
		name = path.Base(t.Abs())
	}
	return
}

func (t *source) Info() target.TemplateInfo {
	return t.info
}
