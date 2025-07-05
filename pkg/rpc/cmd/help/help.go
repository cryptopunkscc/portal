package help

import (
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
)

// Check if given helper contains help printer.
func Check(handler cmd.Handler) bool {
	if handler.Name == Name {
		return true
	}
	for _, h := range handler.Sub {
		if Check(h) {
			return true
		}
	}
	return false
}

// Inject help printer recursively to the given handler and all nested sub-handles.
func Inject(handler *cmd.Handler) {
	for i := range handler.Sub {
		Inject(&handler.Sub[i])
	}
	help := NewHandler(handler)
	handler.AddSub(help)
	if handler.Func == nil {
		handler.Func = help.Func
	}
}

func NewHandler(handler *cmd.Handler) cmd.Handler {
	return cmd.Handler{
		Name: Name, Desc: "Print help.",
		Func: NewFunc(handler),
	}
}

const Name = "help h"

func NewFunc(handler *cmd.Handler) func() Handler { return func() Handler { return Handler{*handler} } }

type Handler struct{ cmd.Handler }

//goland:noinspection GoUnhandledErrorResult
func (h Handler) MarshalCLI() (help string) {
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
		fmt.Fprintln(w, "Commands:")
		fmt.Fprintln(w)
		for _, sub := range h.Sub {
			if sub.Desc == "" {
				fmt.Fprintf(w, "\t%s\t\n", formatName(sub.Names()))
			} else {
				fmt.Fprintf(w, "\t%s\t- %s\n", formatName(sub.Names()), sub.Desc)
			}
		}
		fmt.Fprint(w)
	}
	w.Flush()

	return buffer.String()
}

func formatName(names []string) string { return strings.Join(names, " ") }
