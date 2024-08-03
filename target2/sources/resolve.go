package sources

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/exec"
	golang "github.com/cryptopunkscc/portal/target2/go"
	"github.com/cryptopunkscc/portal/target2/html"
	"github.com/cryptopunkscc/portal/target2/js"
)

func Resolver[A target.Base]() target.Resolve[A] {
	return target.Any[A](
		target.Skip("node_modules"),
		target.Try(golang.ResolveProject),
		target.Try(js.ResolveProject),
		target.Try(js.ResolveDist),
		target.Try(js.ResolveBundle),
		target.Try(html.ResolveProject),
		target.Try(html.ResolveDist),
		target.Try(html.ResolveBundle),
		target.Try(exec.ResolveDist),
		target.Try(exec.ResolveBundle),
	)
}
