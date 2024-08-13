package apps

import (
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/html"
	"github.com/cryptopunkscc/portal/resolve/js"
	"github.com/cryptopunkscc/portal/target"
)

var ResolveAll = Resolver[target.App_]()

func Resolver[T target.App_]() target.Resolve[T] {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(js.ResolveBundle),
		target.Try(js.ResolveDist),
		target.Try(html.ResolveBundle),
		target.Try(html.ResolveDist),
		target.Try(exec.ResolveBundle),
		target.Try(exec.ResolveDist),
	)
}
