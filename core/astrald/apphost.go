package astrald

import (
	"context"
	"fmt"
	libApphost "github.com/cryptopunkscc/astrald/lib/apphost"
	modApphostSrc "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"time"
)

const apphostYaml = "apphost.yaml"

func (r *Initializer) initApphostConfig() {
	if r.apphostConfig == nil {
		if err := r.readApphostConfig(); err != nil {
			r.apphostConfig = &modApphostSrc.Config{}
			return
		}
		r.log.Println("loaded existing apphost config")
	}
}

func (r *Initializer) readApphostConfig() (err error) {
	return r.resources.ReadYaml(apphostYaml, &r.apphostConfig)
}
func (r *Initializer) writeApphostConfig() (err error) {
	return r.resources.WriteYaml(apphostYaml, r.apphostConfig)
}

func (r *Initializer) apphostResolveEndpoint() {
	for _, endpoint := range r.apphostConfig.Listen {
		r.Apphost.Endpoint = endpoint
		return
	}
	r.Apphost.Endpoint = libApphost.DefaultEndpoint
}

func (r *Initializer) apphostSetupAuthToken(pkg string) (err error) {
	plog.TraceErr(&err)
	t, ok := r.ResolvedTokens.Get(pkg)
	if !ok {
		return fmt.Errorf("no token found for %s", pkg)
	}
	r.Apphost.AuthToken = string(t.Token)
	return
}

func (r *Initializer) apphostIsRunning() bool {
	return r.Apphost.Reconnect() == nil
}

func (r *Initializer) apphostAwait(ctx context.Context) (err error) {
	log := plog.Get(ctx).D()
	retry := flow.Await{
		Delay: 50 * time.Millisecond,
		UpTo:  5 * time.Second,
		Mod:   6,
		Ctx:   ctx,
	}.Chan()
	for n := range retry {
		log.Println("awaiting apphost:", n)
		err = r.Apphost.Connect()
		if err == nil || err.Error() == "token authentication failed" {
			err = nil
			log.Println("apphost started")
			return
		}
	}
	return
}
