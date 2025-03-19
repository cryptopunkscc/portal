package main

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/build_all"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/template"
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
			{Type: "string", Desc: "List of templates with optional module names like: 'svelte js' or 'svelte:frontend js:backend'."},
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
	dir = target.Abs(dir)
	parsedTargets := parseTargets(targets)
	if len(parsedTargets) == 0 {
		return errors.New("no targets specified")
	}
	if err = template.NewRunner(dir).GenerateProjects(parsedTargets); err != nil {
		return
	}
	if err = build_all.NewRunner().Dist(ctx, dir); err != nil {
		return
	}
	log.Println("* done")
	return
}

func parseTargets(targets string) (out map[string]string) {
	out = make(map[string]string)
	for _, s := range strings.Split(targets, " ") {
		chunks := strings.Split(s, ":")
		tmpl := chunks[0]
		name := tmpl
		if len(chunks) > 1 {
			name = chunks[1]
		}
		out[name] = tmpl
	}
	return
}
