package template

import (
	"github.com/cryptopunkscc/portal/target"
	"path/filepath"
)

type source struct {
	target.Source
	info target.TemplateInfo
}

func (t *source) Name() (name string) {
	name = t.info.ShortName
	if name == "" {
		name = filepath.Base(t.Abs())
	}
	return
}

func (t *source) Info() target.TemplateInfo {
	return t.info
}
