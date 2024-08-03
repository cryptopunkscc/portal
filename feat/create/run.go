package create

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/template"
	"github.com/cryptopunkscc/portal/target"
)

type (
	Dist    func(context.Context, ...string) error
	Factory func(dir string, templates map[string]string) target.Run[target.Template]
)

type Feat struct {
	factory Factory
	dist    Dist
}

func NewFeat(
	factory Factory,
	dist Dist,
) *Feat {
	return &Feat{
		factory: factory,
		dist:    dist,
	}
}

func (f Feat) Run(
	ctx context.Context,
	dir string,
	targets map[string]string,
) (err error) {
	log := plog.Get(ctx).Type(f).Set(&ctx)
	create := f.factory(dir, targets)
	for _, t := range target.List(
		template.Resolve,
		source.Embed(template.TemplatesFs),
	) {
		if _, ok := targets[t.Name()]; !ok {
			continue
		}
		if err = create(ctx, t); err != nil {
			log.E().Printf("Error creating project from template: %v", err)
		}
	}

	// sanity check
	return f.dist(ctx, dir)
}
