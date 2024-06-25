package template

import (
	"embed"
	targetSource "github.com/cryptopunkscc/portal/target/source"
	"io/fs"
)

//go:embed all:tmpl
var TemplatesFs embed.FS

var CommonsFs fs.FS

func init() {
	if CommonsFs = targetSource.FromFS(TemplatesFs, "tmpl/common").Lift().Files(); CommonsFs == nil {
		panic("cannot resolve templates commons.")
	}
}
