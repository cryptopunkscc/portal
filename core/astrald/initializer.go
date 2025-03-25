package astrald

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	modApphostSrc "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/portal/api/astrald"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/resources"
)

type Initializer struct {
	AgentAlias string
	NodeRoot   mem.String
	TokensDir  mem.String
	Apphost    *apphost.Adapter
	Runner     astrald.Runner

	log            plog.Logger
	resources      resources.FileResources
	nodeIdentity   *astral.Identity
	apphostConfig  *modApphostSrc.Config
	nodeToken      string
	restartAstrald bool
}

func (r *Initializer) Start(ctx context.Context) (err error) {
	r.log = plog.Get(ctx).Type(r)
	if !r.isInitialized() {
		if err = r.initialize(ctx); err != nil {
			return
		}
	}
	err = r.start(ctx)
	return
}

func (r *Initializer) isInitialized() bool {
	return r.fetchAuthToken(r.AgentAlias) == nil
}

func (r *Initializer) initialize(ctx context.Context) (err error) {
	// try to resolve and set apphost endpoint from config.
	if err = r.initNodeResources(); err != nil {
		return
	}
	r.initApphostConfig()
	r.apphostResolveEndpoint()

	// try to resolve node auth token and set to apphost
	if err = r.readOrGenerateNodeIdentity(); err != nil {
		return
	}
	if err = r.resolveNodeAuthToken(); err != nil {
		return
	}
	r.Apphost.AuthToken = r.nodeToken
	if !r.apphostIsRunning() {
		if err = r.startAstrald(ctx); err != nil {
			return
		}
		if err = r.apphostAwait(ctx); err != nil {
			return
		}
	} else if r.restartAstrald {
		return plog.Errorf("cannot configure node token: astrald already running")
	}

	if err = r.removeTemporaryNodeAuthToken(); err != nil {
		return
	}
	if err = r.resolveAuthToken(r.AgentAlias); err != nil {
		return
	}
	return
}

func (r *Initializer) start(ctx context.Context) (err error) {
	r.verifyAgentToken()

	// try to resolve and set apphost endpoint from config.
	if err := r.initNodeResources(); err == nil {
		r.initApphostConfig()
		r.apphostResolveEndpoint()
	}

	// if apphost is not running start astrald and await apphost interface
	if !r.apphostIsRunning() {
		if err = r.startAstrald(ctx); err != nil {
			return
		}
		if err = r.apphostAwait(ctx); err != nil {
			return
		}
	}
	return
}

// verify the agent access token has been set
func (r *Initializer) verifyAgentToken() {
	if len(r.Apphost.AuthToken) == 0 || r.Apphost.AuthToken == r.nodeToken {
		panic(fmt.Errorf("invalid agent token with len %d", len(r.Apphost.AuthToken)))
	}
}
