package dev

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"log"
	"sync"
	"time"
)

type Feat struct {
	wait  *sync.WaitGroup
	spawn target.Spawn
}

func NewFeat(wait *sync.WaitGroup, spawn target.Spawn) target.Spawn {
	return Feat{wait: wait, spawn: spawn}.Run
}

func (f Feat) Run(ctx context.Context, src string) (err error) {

	port := "dev.portal"

	if err = ping(port); err == nil {
		return errors.New("portal dev already running")
	}

	ctx, cancel := context.WithCancel(ctx)
	f.wait.Add(1)
	go func() {
		defer cancel()
		defer f.wait.Done()
		handlers := rpc.Handlers{
			"ping":    func() {},
			"open":    f.spawn,
			"observe": apps.Observe,
		}

		if err = serve.NewRunner(handlers).Run(ctx, port); err != nil {
			log.Printf("serve exit: %v\n", err)
		} else {
			log.Println("serve exit")
		}
	}()

	if err = exec.Retry(5*time.Second, func(i int, i2 int, duration time.Duration) error {
		return ping(port)
	}); err != nil {
		return
	}

	return f.spawn(ctx, src)
}

func ping(port string) error {
	return rpc.Command(rpc.NewRequest(id.Anyone, port), "ping")
}
