package core

import (
	"context"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/ops"
	"github.com/cryptopunkscc/portal/mobile"
	"github.com/cryptopunkscc/portal/pkg/client"
	"github.com/cryptopunkscc/portal/pkg/runner/astrald/debug"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
)

type Service struct {
	ctx     *astral.Context
	cancel  context.CancelFunc
	api     mobile.Api
	astrald debug.Astrald
	client  client.Astrald
	status  int32
}

func (srv *Service) Stop() {
	srv.cancel()
}

func (srv *Service) Start() {
	go srv.start()
}

func (srv *Service) start() {
	srv.set(mobile.STARTING)

	if err := srv.astrald.Start(srv.ctx); err != nil {
		plog.Println(err)
		srv.err(err)
		srv.set(mobile.STOPPED)
		return
	}

	time.Sleep(1 * time.Second)

	srv.client.Init()

	go func() {
		if err := srv.serve(); err != nil {
			srv.err(err)
		}
	}()

	time.Sleep(1 * time.Second)

	srv.set(mobile.STARTED)

	<-srv.ctx.Done()
	err := srv.ctx.Err()
	plog.Println(err)
	srv.err(err)
	srv.set(mobile.STOPPED)
	return
}

func (srv *Service) set(status int32) {
	srv.status = status
	srv.api.Status(status)
}

func (srv *Service) err(err error) {
	if err != nil {
		srv.api.Error(err.Error())
		plog.Println(err)
	}
}

func (srv *Service) serve() (err error) {
	defer plog.TraceErr(&err)
	set := ops.NewSet()
	err = set.AddSubSet("portal", ops.Struct(srv, "Op"))
	if err != nil {
		return
	}
	if err = ops.Serve(srv.ctx, set); err != nil {
		return
	}
	return
}
