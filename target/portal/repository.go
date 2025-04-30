package portal

import "github.com/cryptopunkscc/portal/api/target"

func Repository(sources ...target.Source) *target.SourcesRepository[target.Portal_] {
	return &target.SourcesRepository[target.Portal_]{
		Sources: sources,
		Resolve: Resolve_,
	}
}
