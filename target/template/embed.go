package template

import (
	"embed"
	targetSource "github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
)

//go:embed all:tmpl
var TemplatesFs embed.FS

var CommonsFs fs.FS

func init() {
	if CommonsFs = targetSource.Resolve(TemplatesFs, "tmpl/common").Lift().Files(); CommonsFs == nil {
		panic("cannot resolve templates commons.")
	}
}
