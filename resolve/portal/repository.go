package portal

import "github.com/cryptopunkscc/portal/api/target"

func Repository(sources ...target.Source) *target.SourcesRepository {
	return &target.SourcesRepository{
		Sources: sources,
		Resolve: Resolve_,
	}
}
