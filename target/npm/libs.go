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

var libsDefault []target.NodeModule

func LibsEmbed() []target.NodeModule {
	if len(libsDefault) == 0 {
		libsDefault = resolveLibs.List(source.Embed(js.PortalLibFS))
	}
	return libsDefault
}
