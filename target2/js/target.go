package js

import (
	"errors"
	. "github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/bundle"
	"github.com/cryptopunkscc/portal/target2/dist"
	"github.com/cryptopunkscc/portal/target2/npm"
	"io/fs"
)

var ResolveProject = npm.Resolver[Js](ResolveDist)
var ResolveBundle = bundle.Resolver[Js](ResolveDist)
var ResolveDist = dist.Resolver[Js](ResolveJs)

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
