package cli

import (
	"fmt"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"strings"
)

func injectHelp(handler *cmd.Handler) {
	for i := range handler.Sub {
		injectHelp(&handler.Sub[i])
	}
	help := newHelpHandler(handler)
	handler.AddSub(help)
	if handler.Func == nil {
		handler.Func = help.Func
	}
}

func newHelpHandler(handler *cmd.Handler) cmd.Handler {
	return cmd.Handler{
		Name: "help h", Desc: "Print help.",
		Func: newHelpFunc(handler),
	}
}

func newHelpFunc(handler *cmd.Handler) func() Help { return func() Help { return Help{*handler} } }

type Help struct{ cmd.Handler }

func (h Help) MarshalCLI() (help string) {
	help += strings.Join(h.Names(), ", ")
	if h.Desc != "" {
		help += " - " + h.Desc
	}

	help += "\n\n"
	if len(h.Params) > 0 {
		help += "Parameters\n"
		for _, p := range h.Params {
			n := ""
			if p.Name != "" {
				n = "-" + p.Name
			}
			help += fmt.Sprintf("\t%s %s - %s\n", n, p.Type, p.Desc)
		}
		help += "\n"
	}
	if len(h.Sub) > 0 {
		help += "Commands\n"
		for _, sub := range h.Sub {
			help += "\t" + strings.Join(sub.Names(), ", ")
			if sub.Desc != "" {
				help += " - " + sub.Desc
			}
			help += "\n"
		}
		help += "\n"
	}
	return
}
