package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	find "github.com/cryptopunkscc/portal/factory/find/portal"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"log"
	"os"
	"reflect"
	"strings"
)

func main() {
	ctx := context.Background()
	plog.New().D().Set(&ctx)

	err := cli.New(cmd.Handler{
		Name: "portal-list",
		Desc: "Print all targets in given directory.",
		Params: cmd.Params{
			{Type: "string", Desc: "Directory containing targets."},
		},
		Func: listPortals(find.Create[target.Portal_]()),
	}).Run(ctx)

	//cli := clir.NewCli(ctx,
	//	"portal-build",
	//	"Builds portal project and generates application bundle.",
	//	version.Run)
	//cli.Portals(find.Create[target.Portal_]())
	//err := cli.Run()

	if err != nil {
		panic(err)
	}
}

func listPortals(find target.Find[target.Portal_]) func(context.Context, string) error {
	return func(ctx context.Context, path string) (err error) {
		wd, _ := os.Getwd()
		portals, err := find(ctx, path)
		if err != nil {
			return
		}
		for _, source := range portals {
			log.Println(reflect.TypeOf(source), "\t", strings.TrimPrefix(source.Abs(), wd+"/"))
		}
		return
	}
}
