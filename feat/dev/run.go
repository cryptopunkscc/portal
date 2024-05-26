package dev

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"github.com/cryptopunkscc/go-astral-js/target"
	"sync"
	"time"
)

type Feat struct {
	port  string
	wait  *sync.WaitGroup
	query target.Dispatch
	serve target.Dispatch
}

func NewFeat(
	port string,
	wait *sync.WaitGroup,
	spawn target.Dispatch,
	query target.Dispatch,
) target.Dispatch {
	handlers := rpc.Handlers{
		"ping":    func() {},
		"open":    spawn,
		"observe": apps.Observe,
	}
	return Feat{
		port:  port,
		wait:  wait,
		query: query,
		serve: serve.NewRunner(handlers).Run,
	}.Run
}

func (f Feat) Run(ctx context.Context, src string, args ...string) (err error) {
	log := plog.Get(ctx).Type(f).Set(&ctx)
	if err = ping(f.port); err == nil {
		return errors.New("portal dev already running")
	}
	ctx, cancel := context.WithCancel(ctx)
	f.wait.Add(1)
	go func() {
		defer cancel()
		defer f.wait.Done()

		if err = f.serve(ctx, f.port); err != nil {
			log.Printf("serve exit: %v", err)
		} else {
			log.Println("serve exit")
		}
	}()
	if err = exec.Retry(ctx, 5*time.Second, func(i int, i2 int, duration time.Duration) error {
		return ping(f.port)
	}); err != nil {
		return
	}
	return f.query(ctx, src, args...)
}

func ping(port string) error {
	return rpc.Command(rpc.NewRequest(id.Anyone, port), "ping")
}
