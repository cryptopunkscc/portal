package template

import (
	"embed"
	targetSource "github.com/cryptopunkscc/portal/resolve/source"
	"io/fs"
)

//go:embed all:tmpl
var TemplatesFs embed.FS

var CommonsFs fs.FS

func init() {
	sub, err := targetSource.Embed(TemplatesFs).Sub("tmpl/common")
	if err != nil {
		return
	}
	if CommonsFs = sub.Files(); CommonsFs == nil {
		panic("cannot resolve templates commons.")
	}
}
