package apps

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/exec"
	"github.com/cryptopunkscc/portal/target2/html"
	"github.com/cryptopunkscc/portal/target2/js"
)

var ResolveAll target.Resolve[target.App_] = Resolver[target.App_]()

func Resolver[T target.App_]() func(target.Source) (T, error) {
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
