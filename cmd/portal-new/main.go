package main

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/factory/create"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/template"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"log"
	"strings"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)

	err := cli.New(cmd.Handler{
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
	}).Run(ctx)

	//cli := clir.NewCli(ctx,
	//	"portal-new",
	//	"Create new portal project from template.",
	//	version.Run)
	//cli.Create(
	//	template.List,
	//	create.Create())
	//err := cli.Run()

	if err != nil {
		panic(err)
	}
	log.Println("* done")
}

func createProject(ctx context.Context, targets string, dir string) error {
	if targets == "" {
		return errors.New("no targets specified")
	}
	if dir == "" {
		dir = "."
	}
	parsedTargets := parseTargets(targets)
	return create.Run(ctx, dir, parsedTargets)
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
