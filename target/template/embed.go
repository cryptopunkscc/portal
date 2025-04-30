package template

import (
	"embed"
	"github.com/cryptopunkscc/portal/target/source"
	"io/fs"
)

//go:embed all:tmpl
var TemplatesFs embed.FS

var CommonsFs fs.FS

func init() {
	sub, err := source.Embed(TemplatesFs).Sub("tmpl/common")
	if err != nil {
		return
	}
	if CommonsFs = sub.FS(); CommonsFs == nil {
		panic("cannot resolve templates commons.")
	}
}
