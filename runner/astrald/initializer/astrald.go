package initializer

import (
	"context"
	"fmt"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/resources"
	"github.com/cryptopunkscc/portal/runner/astrald"
)

type Astrald struct {
	AgentAlias string
	NodeRoot   string
	TokensDir  string
	Config     astrald.Config
	Runner     astrald.Runner
	Apphost    *apphost.Adapter

	log            plog.Logger
	resources      resources.Dir
	nodeIdentity   *astral.Identity
	nodeToken      string
	restartAstrald bool
}

func (i *Astrald) Start(ctx context.Context) (err error) {
	i.log = plog.Get(ctx).Type(i)
	if !i.isInitialized() {
		if err = i.initialize(ctx); err != nil {
			return
		}
	}
	err = i.start(ctx)
	return
}

func (i *Astrald) isInitialized() bool {
	return i.fetchAuthToken(i.AgentAlias) == nil
}

func (i *Astrald) initialize(ctx context.Context) (err error) {
	// try to resolve and set apphost endpoint from config.
	if err = i.initNodeResources(); err != nil {
		return
	}
	if err = i.createConfigs(); err != nil {
		return
	}

	i.initApphostConfig()
	i.apphostResolveEndpoint()

	// try to resolve node auth token and set to apphost
	if err = i.readOrGenerateNodeIdentity(); err != nil {
		return
	}
	if err = i.resolveNodeAuthToken(); err != nil {
		return
	}
	i.Apphost.Token = i.nodeToken
	if !i.apphostIsRunning() {
		if err = i.startAstrald(ctx); err != nil {
			return
		}
		if err = i.apphostAwait(ctx); err != nil {
			return
		}
	} else if i.restartAstrald {
		return plog.Errorf("cannot configure node token: astrald already running")
	}

	if err = i.removeTemporaryNodeAuthToken(); err != nil {
		return
	}
	if err = i.resolveAuthToken(i.AgentAlias); err != nil {
		return
	}
	return
}

func (i *Astrald) start(ctx context.Context) (err error) {
	i.verifyAgentToken()

	// try to resolve and set apphost endpoint from config.
	if err := i.initNodeResources(); err == nil {
		i.initApphostConfig()
		i.apphostResolveEndpoint()
	}

	// if apphost is not running start astrald and await apphost interface
	if !i.apphostIsRunning() {
		if err = i.startAstrald(ctx); err != nil {
			return
		}
		if err = i.apphostAwait(ctx); err != nil {
			return
		}
	}
	return
}

// verify the agent access token has been set
func (i *Astrald) verifyAgentToken() {
	if len(i.Apphost.Token) == 0 || i.Apphost.Token == i.nodeToken {
		panic(fmt.Errorf("invalid agent token with len %d", len(i.Apphost.Token)))
	}
}
