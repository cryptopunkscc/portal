package create

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/target/template"
)

type Feat struct {
	newCreate func(dir string, templates map[string]string) func(Template) error
	dist      func(context.Context, ...string) error
}

func NewFeat(
	newCreate func(dir string, templates map[string]string) func(Template) error,
	dist func(context.Context, ...string) error,
) *Feat {
	return &Feat{
		newCreate: newCreate,
		dist:      dist,
	}
}

func (f Feat) Run(
	ctx context.Context,
	dir string,
	targets map[string]string,
) (err error) {
	log := plog.Get(ctx).Type(f).Set(&ctx)
	create := f.newCreate(dir, targets)
	resolve := Any[Template](Try(template.Resolve))
	src := source.FromFS(template.TemplatesFs)

	for _, t := range source.List(resolve, src) {
		if _, ok := targets[t.Name()]; !ok {
			continue
		}
		if err = create(t); err != nil {
			log.E().Printf("Error creating project from template: %v", err)
		}
	}

	// sanity check
	return f.dist(ctx, dir)
}
