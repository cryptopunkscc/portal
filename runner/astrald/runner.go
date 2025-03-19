package astrald

import (
	"context"
	modApphostSrc "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/resources"
	"github.com/labstack/gommon/random"
	"time"
)

type Runner struct {
	Astrald
	exec.Cmd
}

type Astrald interface {
	RootDir() string
	Start(ctx context.Context) error
}

func (r *Runner) Start(ctx context.Context) (err error) {
	defer plog.TraceErr(&err)
	plog.Get(ctx).Type(r).Set(&ctx)
	if err = r.Astrald.Start(ctx); err != nil {
		return
	}
	err = r.awaitApphost(ctx)
	return
}

func (r *Runner) awaitApphost(ctx context.Context) (err error) {
	log := plog.Get(ctx).D()
	adapter := apphost.Adapter{}
	if r.RootDir() != "" {
		adapter.Endpoint = r.apphostAddr()
		adapter.AuthToken = random.String(10) // fake access token is enough to simulate a ping.
	}
	retry := flow.Await{
		Delay: 50 * time.Millisecond,
		UpTo:  5 * time.Second,
		Mod:   6,
		Ctx:   ctx,
	}.Chan()
	for n := range retry {
		log.Println("awaiting apphost:", n)
		err = adapter.Connect()
		if err == nil || err.Error() == "token authentication failed" {
			err = nil
			log.Println("apphost started")
			return
		}
	}
	return
}

func (r *Runner) apphostAddr() (address string) {
	res, err := resources.NewFileResources(r.RootDir())
	if err != nil {
		return
	}
	config := modApphostSrc.Config{}
	if err = res.ReadYaml("apphost.yaml", &config); err != nil {
		return
	}
	for _, s := range config.Listen[:1] {
		address = s
	}
	return
}
