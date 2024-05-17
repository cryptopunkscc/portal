package dev

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"log"
)

type Feat struct {
	spawn target.Spawn
}

func NewFeat(spawn target.Spawn) target.Spawn {
	return Feat{spawn: spawn}.Run
}

func (f Feat) Run(ctx context.Context, src string) (err error) {

	port := "dev.portal"

	if err = rpc.Command(rpc.NewRequest(id.Anyone, port), "ping"); err != nil {
		return errors.New("portal dev already running")
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		handlers := rpc.Handlers{
			"ping":    func() {},
			"open":    f.spawn,
			"observe": apps.Observe,
		}

		if err = serve.Run(ctx, port, handlers); err != nil {
			log.Printf("serve exit: %v\n", err)
		} else {
			log.Println("serve exit")
		}
	}()

	return f.spawn(ctx, src)
}
