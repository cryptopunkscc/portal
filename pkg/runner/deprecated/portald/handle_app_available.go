package portald

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/source/app"
	"github.com/cryptopunkscc/portal/pkg/util/flow"
)

func (s *Service) AvailableApps(ctx context.Context, follow bool) (out flow.Input[app.ReleaseInfo]) {
	return s.AppObjects().Scan(ctx, follow)
}
