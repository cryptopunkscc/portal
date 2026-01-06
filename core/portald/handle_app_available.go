package portald

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/source/app"
)

func (s *Service) AvailableApps(ctx context.Context, follow bool) (out flow.Input[app.ReleaseInfo]) {
	return s.AppObjects().Scan(ctx, follow)
}
