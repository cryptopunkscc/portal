package html

import (
	"errors"
	. "github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/bundle"
	"github.com/cryptopunkscc/portal/target2/dist"
	"github.com/cryptopunkscc/portal/target2/npm"
	"io/fs"
)

var ResolveProject = npm.Resolver[Html](ResolveDist)
var ResolveBundle = bundle.Resolver[Html](ResolveDist)
var ResolveDist = dist.Resolver[Html](ResolveHtml)

type html struct{ Source }

func (h html) IndexHtml() {}

func ResolveHtml(src Source) (t Html, err error) {
	stat, err := fs.Stat(src.Files(), "index.html")
	if err != nil {
		return
	}
	if stat.IsDir() {
		return nil, errors.New("index.html is not a file")
	}
	return html{Source: src}, nil
}
