package html

import (
	"errors"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/dist"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"io/fs"
)

type html struct{ Source }

func (h html) IndexHtml() {}

func ResolveHtml(src Source) (t Html, err error) {
	stat, err := fs.Stat(src.FS(), "index.html")
	if err != nil {
		return
	}
	if stat.IsDir() {
		return nil, errors.New("index.html is not a file")
	}
	return html{Source: src}, nil
}

var ResolveDist = dist.Resolver[Html](ResolveHtml)
var ResolveBundle = bundle.Resolver[Html](ResolveDist)
var ResolveProject = npm.Resolver[Html](ResolveHtml)
