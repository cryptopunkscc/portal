package sources

import (
	"github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/exec"
	golang "github.com/cryptopunkscc/portal/target2/go"
	"github.com/cryptopunkscc/portal/target2/html"
	"github.com/cryptopunkscc/portal/target2/js"
)

var ResolveAll = Resolver[target2.Base]()

func Resolve[A target2.Base](src target2.Source) (A, error) {
	return Resolver[A]()(src)
}

func Resolver[A target2.Base]() target2.Resolve[A] {
	return target2.Any[A](
		target2.Skip("node_modules"),
		target2.Try(golang.ResolveProject),
		target2.Try(js.ResolveProject),
		target2.Try(js.ResolveDist),
		target2.Try(js.ResolveBundle),
		target2.Try(html.ResolveProject),
		target2.Try(html.ResolveDist),
		target2.Try(html.ResolveBundle),
		target2.Try(exec.ResolveDist),
		target2.Try(exec.ResolveBundle),
	)
}
