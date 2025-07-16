package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/target/bundle"
)

func (s *Service) AvailableApps(ctx context.Context, follow bool) (out flow.Input[bundle.Info], err error) {
	return s.Bundles().Scan(ctx, follow)
}
