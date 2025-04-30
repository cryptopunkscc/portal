package html

import (
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/dist"
	"github.com/cryptopunkscc/portal/resolve/npm"
	"io/fs"
)

var ResolveDist = dist.Resolver[target.Html](ResolveHtml)
var ResolveBundle = bundle.Resolver[target.Html](ResolveDist)
var ResolveProject = npm.Resolver[target.Html](ResolveHtml)

func ResolveHtml(source target.Source) (html target.Html, err error) {
	defer plog.TraceErr(&err)
	s, err := fs.Stat(source.FS(), "index.html")
	if err != nil {
		return
	}
	if s.IsDir() {
		return nil, errors.New("index.html is not a file")
	}
	html = Source{source}
	return
}
