package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/portal"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	_ "github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/dir"
	"github.com/cryptopunkscc/portal/runtime/rpc/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"

	"log"
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
				Func: listPortals(find[target.Portal_]()),
			},
		},
	}).Run(ctx)
	if err != nil {
		panic(err)
	}
}

func find[T target.Portal_]() target.Find[T] {
	return target.FindByPath(
		source.File,
		sources.Resolver[T]()).
		OrById(path.Resolver(portal.Resolve_, dir.AppSource))
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
