package template

import (
	"github.com/cryptopunkscc/portal/api/target"
	"path/filepath"
)

type Source struct {
	target.Source
	info target.TemplateInfo
}

func (t *Source) Info() target.TemplateInfo { return t.info }
func (t *Source) Name() (name string) {
	name = t.info.ShortName
	if name == "" {
		name = filepath.Base(t.Abs())
	}
	return
}
