package sources

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/exec"
	golang "github.com/cryptopunkscc/portal/resolve/go"
	"github.com/cryptopunkscc/portal/resolve/html"
	"github.com/cryptopunkscc/portal/resolve/js"
)

func Resolver[A target.Portal_]() target.Resolve[A] {
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
