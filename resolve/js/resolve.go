package js

import (
	"errors"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/dist"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"io/fs"
)

type js struct{}

func (h js) MainJs() {}

func ResolveJs(src Source) (t Js, err error) {
	stat, err := fs.Stat(src.Files(), "main.js")
	if err != nil {
		return
	}
	if stat.IsDir() {
		return nil, errors.New("main.js is not a file")
	}
	t = js{}
	return
}

var ResolveDist = dist.Resolver[Js](ResolveJs)
var ResolveBundle = bundle.Resolver[Js](ResolveDist)
var ResolveProject = npm.Resolver[Js](ResolveJs)
