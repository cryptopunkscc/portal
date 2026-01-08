package main

import (
	"fmt"
	"strings"

	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/os"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func main() { cli.Run(handler) }

var handler = cmd.Handler{
	Name: "portal-tools",
	Sub: cmd.Handlers{
		cmd.Handler{
			Func: listGoImports,
			Name: "li",
			Desc: "List go imports.",
			Params: cmd.Params{
				cmd.Param{Type: "string", Desc: "Path to go file."},
			},
		},
	},
}

func listGoImports(path string) (err error) {
	imports, err := golang.ListImports(os.Abs(path))
	if err != nil {
		return err
	}
	for i, s := range imports {
		print(fmt.Sprintf("%d %s (%d) ", i, s.Import, len(s.Refs)))
		//if len(s.Refs) < 2 {
		//	println(fmt.Sprintf("\t %v", s.Refs))
		//} else {
		//}
		println(":\n\t" + strings.Join(s.Refs, "\n\t"))
	}
	return
}
