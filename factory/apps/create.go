package apps

import (
	"github.com/cryptopunkscc/portal/api/apps"
	"github.com/cryptopunkscc/portal/api/target"
	resolve "github.com/cryptopunkscc/portal/resolve/apps"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	runtime "github.com/cryptopunkscc/portal/runtime/apps"
	"os"
)

func Default() apps.Apps {
	return Dir(runtime.DefaultDir())
}

func Dir(dir string) apps.Apps {
	var err error
	a := runtime.Apps{}
	a.File = source.File
	a.Resolve = resolve.Resolver[target.App_]()
	a.Find = target.FindByPath(
		source.File,
		sources.Resolver[target.Bundle_](),
	)
	if err = os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	if a.Dir, err = a.File(dir); err != nil {
		panic(err)
	}
	return a
}
