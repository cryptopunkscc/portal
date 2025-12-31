package main

import (
	"fmt"
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/cryptopunkscc/portal/source/html"
	"github.com/cryptopunkscc/portal/source/js"
)

func init() { cmd.DefaultHandlers.Add(ListTargetsHandler) }

var ListTargetsHandler = cmd.Handler{
	Name: "lt",
	Desc: "List apps and projects recursively found in given path.",
	Func: ListTargets,
}

type ListImportsOpt struct {
	Local bool `cli:"local l"`
}

func ListTargets(src string) (err error) {
	s := source.Providers{
		source.OsFs,
	}.GetSource(src)
	if s == nil {
		return fs.ErrNotExist
	}
	for i, ss := range source.CollectT[app.App](s,
		&html.App{},
		&html.Bundle{},
		&html.Project{},
		&js.App{},
		&js.Bundle{},
		&js.Project{},
	) {
		println(fmt.Sprintf("%d. %T:%s %v", i, ss, ss.GetPath(), ss.GetMetadata().Manifest))
	}
	return
}
