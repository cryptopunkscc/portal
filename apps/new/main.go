package main

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/factory/create"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/template"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"log"
	"strings"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: createProject,
		Name: "portal-new",
		Desc: "Create new portal project from template.",
		Params: cmd.Params{
			{Type: "string", Desc: "List of templates with optional module names like: 'svelte backend' or 'svelte:front backend:back'."},
			{Type: "string", Desc: "Project directory."},
		},
		Sub: cmd.Handlers{
			{
				Func: template.List,
				Name: "list l",
				Desc: "List available templates.",
			},
		},
	}
}

func createProject(ctx context.Context, targets string, dir string) (err error) {
	if targets == "" {
		return errors.New("no targets specified")
	}
	if dir == "" {
		dir = "."
	}
	parsedTargets := parseTargets(targets)
	err = create.Run(ctx, dir, parsedTargets)
	if err != nil {
		log.Println("* done")
	}
	return
}

func parseTargets(targets string) (out map[string]string) {
	out = make(map[string]string)
	for _, s := range strings.Split(targets, " ") {
		chunks := strings.Split(s, ":")
		template := chunks[0]
		name := template
		if len(chunks) > 1 {
			name = chunks[1]
		}
		out[template] = name
	}
	return
}
