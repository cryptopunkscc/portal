package js

import (
	"embed"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
	"os"
)

//go:embed portal
var PortalLibFS embed.FS

var resolve = target.Any[target.NodeModule](
	target.Skip("node_modules"),
	target.Try(npm.ResolveNodeModule),
)

var LibsDefault = LibsEmbed()

func LibsEmbed() []target.NodeModule {
	return resolve.List(source.Embed(PortalLibFS))
}

func LibsDir() []target.NodeModule {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("could not get wd: %w", err))
	}
	root, err := golang.FindProjectRoot(wd)
	if err != nil {
		panic(fmt.Errorf("could not find project root: %w", err))
	}
	return resolve.List(source.Dir(root, "core", "js", "embed", "portal"))
}
