package create

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/runner/create"
	. "github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"github.com/cryptopunkscc/go-astral-js/target/template"
)

type Feat struct{}

func NewFeat() *Feat { return &Feat{} }

func (f Feat) Run(
	ctx context.Context,
	dir string,
	targets map[string]string,
) (err error) {
	log := plog.Get(ctx).Type(f).Set(&ctx)
	runner := create.NewRunner(dir, targets)
	resolve := Any[Template](Try(template.Resolve))
	src := source.FromFS(template.TemplatesFs)

	for _, t := range source.List(resolve, src) {
		if _, ok := targets[t.Name()]; !ok {
			continue
		}
		if err = runner.Run(t); err != nil {
			log.E().Printf("Error creating project from template: %v", err)
		}
	}

	// sanity check
	return build.NewFeat().Dist(ctx, dir)
}
