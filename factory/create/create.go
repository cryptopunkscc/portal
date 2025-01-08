package create

import (
	"github.com/cryptopunkscc/portal/factory/build"
	"github.com/cryptopunkscc/portal/runner/create"
	"github.com/cryptopunkscc/portal/runner/template"
)

func Create() create.Run {
	return create.Runner(
		template.Runner,
		build.Create().Dist,
	)
}

var Run = create.Runner(
	template.Runner,
	build.Create().Dist,
)
