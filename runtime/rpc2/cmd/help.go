package cmd

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
)

func HasHelp(handler Handler) bool {
	if handler.Name == HelpName {
		return true
	}
	for _, h := range handler.Sub {
		if HasHelp(h) {
			return true
		}
	}
	return false
}

func InjectHelp(handler *Handler) {
	for i := range handler.Sub {
		InjectHelp(&handler.Sub[i])
	}
	help := NewHelpHandler(handler)
	handler.AddSub(help)
	if handler.Func == nil {
		handler.Func = help.Func
	}
}

func NewHelpHandler(handler *Handler) Handler {
	return Handler{
		Name: HelpName, Desc: "Print help.",
		Func: NewHelpFunc(handler),
	}
}

const HelpName = "help h"

func NewHelpFunc(handler *Handler) func() Help { return func() Help { return Help{*handler} } }

type Help struct{ Handler }

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
