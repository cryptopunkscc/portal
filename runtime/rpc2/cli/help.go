package cli

import (
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
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

//goland:noinspection GoUnhandledErrorResult
func (h Help) MarshalCLI() (help string) {
	buffer := &bytes.Buffer{}
	w := tabwriter.NewWriter(buffer, 4, 4, 2, ' ', 0)

	name := formatName(h.Names())
	if h.Desc != "" {
		name += " - " + h.Desc
	}

	fmt.Fprintln(w, name)
	fmt.Fprintln(w)

	if len(h.Params) > 0 {
		fmt.Fprintln(w, "Parameters:")
		fmt.Fprintln(w)
		i := 0
		for _, p := range h.Params {
			n := ""
			if p.Name != "" {
				fields := strings.Fields(p.Name)
				slices.Reverse(fields)
				for i, f := range fields {
					fields[i] = "-" + f
				}
				n += strings.Join(fields, " ")
			} else {
				n = "$" + strconv.Itoa(i)
				i++
			}
			fmt.Fprintf(w, "\t%s\t[%s]\t- %s\n", n, p.Type, p.Desc)
		}
		fmt.Fprintln(w)
	}
	if len(h.Sub) > 0 {
		fmt.Fprintln(w, "Subcommands:")
		fmt.Fprintln(w)
		for _, sub := range h.Sub {
			fmt.Fprintf(w, "\t%s\t- %s\n", formatName(sub.Names()), sub.Desc)
		}
		fmt.Fprintln(w)
	}
	w.Flush()

	return buffer.String()
}

func formatName(names []string) string { return strings.Join(names, " ") }
