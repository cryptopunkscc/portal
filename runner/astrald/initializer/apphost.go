package initializer

import (
	"context"
	"time"

	libApphost "github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

const apphostYaml = "apphost.yaml"

func (i *Astrald) initApphostConfig() {
	if err := i.readApphostConfig(); err == nil {
		i.log.Println("loaded existing apphost config")
	}
}

func (i *Astrald) readApphostConfig() (err error) {
	return i.resources.ReadYaml(apphostYaml, &i.Config.Apphost)
}
func (i *Astrald) writeApphostConfig() (err error) {
	return i.resources.WriteYaml(apphostYaml, i.Config.Apphost)
}

func (i *Astrald) apphostResolveEndpoint() {
	for _, endpoint := range i.Config.Apphost.Listen {
		i.Apphost.Endpoint = endpoint
		return
	}
	i.Apphost.Endpoint = libApphost.DefaultEndpoint
}

func (i *Astrald) apphostIsRunning() bool {
	return i.Apphost.Reconnect() == nil
}

func (i *Astrald) apphostAwait(ctx context.Context) (err error) {
	log := plog.Get(ctx).D()
	retry := flow.Await{
		Delay: 50 * time.Millisecond,
		UpTo:  5 * time.Second,
		Mod:   6,
		Ctx:   ctx,
	}.Chan()
	for n := range retry {
		log.Println("awaiting apphost:", n)
		err = i.Apphost.Connect()
		if err == nil || err.Error() == "token authentication failed" {
			err = nil
			log.Println("apphost started")
			return
		}
	}
	return
}
