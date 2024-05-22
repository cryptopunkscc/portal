package create

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/build"
	"github.com/cryptopunkscc/go-astral-js/runner/create"
	. "github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"github.com/cryptopunkscc/go-astral-js/target/template"
	"log"
)

type Feat struct{}

func NewFeat() *Feat { return &Feat{} }

func (f Feat) Run(
	ctx context.Context,
	dir string,
	targets map[string]string,
) (err error) {
	runner := create.NewRunner(dir, targets)
	resolve := Any[Template](Try(template.Resolve))
	src := source.FromFS(template.TemplatesFs)

	for t := range source.Stream(resolve, src) {
		if _, ok := targets[t.Name()]; !ok {
			continue
		}
		if err = runner.Run(t); err != nil {
			log.Printf("Error creating project from template: %v", err)
		}
	}

	// sanity check
	return build.NewFeat().Dist(ctx, dir)
}
