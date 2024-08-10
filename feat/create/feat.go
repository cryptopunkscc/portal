package create

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/template"
	"github.com/cryptopunkscc/portal/target"
)

type (
	Run     func(ctx context.Context, dir string, targets map[string]string) (err error)
	Dist    func(context.Context, ...string) error
	Factory func(dir string, templates map[string]string) target.Run[target.Template]
)

func Feat(factory Factory, dist Dist) Run {
	return func(ctx context.Context, dir string, targets map[string]string) (err error) {
		log := plog.Get(ctx)
		create := factory(dir, targets)
		for _, t := range template.Resolve.List(
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
		return dist(ctx, dir)
	}
}
