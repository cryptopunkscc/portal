package create

import (
	"github.com/cryptopunkscc/portal/factory/build"
	"github.com/cryptopunkscc/portal/feat/create"
	"github.com/cryptopunkscc/portal/runner/template"
)

func Create() *create.Feat {
	return create.NewFeat(
		template.NewRun,
		build.Create().Dist,
	)
}
