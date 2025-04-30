package js

import (
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/dist"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"io/fs"
)

var ResolveDist = dist.Resolver[target.Js](ResolveJs)
var ResolveBundle = bundle.Resolver[target.Js](ResolveDist)
var ResolveProject = npm.Resolver[target.Js](ResolveJs)

func ResolveJs(source target.Source) (js target.Js, err error) {
	defer plog.TraceErr(&err)
	stat, err := fs.Stat(source.FS(), "main.js")
	if err != nil {
		return
	}
	if stat.IsDir() {
		return nil, errors.New("main.js is not a file")
	}
	js = Source{source}
	return
}
