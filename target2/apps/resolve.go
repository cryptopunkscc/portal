package apps

import (
	"github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/exec"
	"github.com/cryptopunkscc/portal/target2/html"
	"github.com/cryptopunkscc/portal/target2/js"
)

var ResolveAll target2.Resolve[target2.Base] = Resolver[target2.Base]()

func Resolve[T target2.Base](src target2.Source) (T, error) { return Resolver[T]()(src) }
func Resolver[T target2.Base]() func(target2.Source) (T, error) {
	return target2.Any[T](
		target2.Skip("node_modules"),
		target2.Try(js.ResolveBundle),
		target2.Try(js.ResolveDist),
		target2.Try(html.ResolveBundle),
		target2.Try(html.ResolveDist),
		target2.Try(exec.ResolveBundle),
		target2.Try(exec.ResolveDist),
	)
}
