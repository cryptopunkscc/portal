package create

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/template"
)

type (
	Run     func(ctx context.Context, dir string, targets map[string]string) (err error)
	Dist    func(context.Context, ...string) error
	Factory func(dir string, templates map[string]string) target.Run[target.Template]
)

func Runner(factory Factory, dist Dist) Run {
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
		if err = dist(ctx, dir); errors.Is(err, target.ErrNotFound) {
			err = nil
		}
		return
	}
}
