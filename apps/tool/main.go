package main

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/resolve/portal"
	"github.com/cryptopunkscc/portal/resolve/source"
	"os"
	"reflect"
	"strings"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)

	err := cli.New(cmd.Handler{
		Name: "tool",
		Desc: "Various portal tools.",
		Sub: cmd.Handlers{
			{
				Name: "list l",
				Desc: "Print all portal targets in given directory.",
				Params: cmd.Params{
					{Type: "string", Desc: "Directory containing targets."},
				},
				Func: list,
			},
		},
	}).Run(ctx)
	if err != nil {
		panic(err)
	}
}

func list(src string) (p []printable) {
	wd, _ := os.Getwd()
	for _, r := range targets.All(src) {
		p = append(p, printable{wd, r})
	}
	return
}

var targets = target.Provider[target.Portal_]{
	Repository: source.Repository,
	Resolve:    target.Any[target.Portal_](portal.Resolve_.Try),
}

type printable struct {
	wd string
	target.Portal_
}

func (p printable) MarshalCLI() string {
	return fmt.Sprintln(reflect.TypeOf(p), "\t", strings.TrimPrefix(p.Abs(), p.wd+"/"))
}
