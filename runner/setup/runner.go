package setup

import (
	"context"
	"github.com/cryptopunkscc/astrald/astral"
	modApphost "github.com/cryptopunkscc/astrald/mod/apphost"
	modApphostSrc "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/resources"
	"github.com/cryptopunkscc/portal/runner/astrald"
	"github.com/cryptopunkscc/portal/runtime/token"
)

type Runner struct {
	NodeRoot     string
	Runner       astrald.Runner
	Tokens       token.Repository
	CreateTokens []string

	ApphostConfig  *modApphostSrc.Config
	ResolvedTokens sig.Map[string, *modApphost.AccessToken]

	// private
	nodeAuthToken string
	log           plog.Logger
	resources     resources.FileResources
	nodeIdentity  *astral.Identity
}

func (r *Runner) Setup(ctx context.Context) (err error) {
	r.log = plog.Get(ctx).Type(r).Set(&ctx)
	if err = r.setupResources(); err != nil {
		return
	}
	if err = r.readOrGenerateNodeIdentity(); err != nil {
		return
	}
	if err = r.resolveNodeAuthToken(); err != nil {
		return
	}
	if err = r.startAstrald(ctx); err != nil {
		return
	}
	if err = r.resolveAuthToken("portald"); err != nil {
		return
	}
	if err = r.removeTemporaryNodeAuthToken(); err != nil {
		return
	}
	return
}
