package main

import (
	"fmt"
	"strings"

	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/os"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func init() { cmd.DefaultHandlers.Add(ListGoImportsHandler) }

var ListGoImportsHandler = cmd.Handler{
	Func: ListGoImports,
	Name: "li",
	Desc: "List go imports.",
	Params: cmd.Params{
		cmd.Param{Type: "string", Desc: "Path to go file."},
	},
}

func ListGoImports(opt ListImportsOpt, path string) (err error) {
	path = os.Abs(path)
	require := ""
	if opt.Local {
		if project, err := golang.ResolveProject(path); err == nil && project != nil {
			require = project.Name
		}
	}
	imports, err := golang.ListImports(path)
	if err != nil {
		return err
	}
	for i, s := range imports {
		if !strings.Contains(s.Import, require) {
			continue
		}
		print(fmt.Sprintf("%d %s (%d) ", i, s.Import, len(s.Refs)))
		println(":\n\t" + strings.Join(s.Refs, "\n\t"))
	}
	return
}
