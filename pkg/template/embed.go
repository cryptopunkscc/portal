package template

import (
	"embed"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/source"
	"io/fs"
)

//go:embed all:tmpl
var TemplatesFs embed.FS

var CommonsFs fs.FS

func init() {
	if CommonsFs = source.Resolve(TemplatesFs, "tmpl/common").Lift().Files(); CommonsFs == nil {
		panic("cannot resolve templates commons.")
	}
}
