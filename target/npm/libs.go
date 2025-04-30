package npm

import (
	"github.com/cryptopunkscc/portal/api/target"
	js "github.com/cryptopunkscc/portal/core/js/embed"
	"github.com/cryptopunkscc/portal/target/source"
)

var resolveLibs = target.Any[target.NodeModule](
	target.Skip("node_modules"),
	target.Try(ResolveNodeModule),
)

var LibsDefault = LibsEmbed()

func LibsEmbed() []target.NodeModule {
	return resolveLibs.List(source.Embed(js.PortalLibFS))
}
