package main

import (
	"fmt"
	"io/fs"

	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/cryptopunkscc/portal/source/html"
	"github.com/cryptopunkscc/portal/source/js"
)

func listTargets(src string) (err error) {
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
