package main

import (
	"context"
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"os/exec"
	"time"
)

func startPortald(ctx context.Context, client portald.Client) (err error) {
	if err = startPortaldProcess(ctx); err != nil {
		return
	}
	if err = awaitPortaldService(ctx, client); err != nil {
		return
	}
	return
}

func startPortaldProcess(ctx context.Context) (err error) {
	plog.Get(ctx).Println("starting portald")
	c := exec.Command("portald")
	err = c.Start()
	return
}

func awaitPortaldService(ctx context.Context, client portald.Client) error {
	log := plog.Get(ctx)
	return flow.Retry(ctx, 8*time.Second, func(i int, n int, d time.Duration) (err error) {
		log.Printf("%d/%d attempt %v: retry after %v", i+1, n, err, d)
		if err = apphost.Default.Connect(); err != nil {
			log.Printf("failed to connect to apphost: %v", err)
			return
		}
		return client.Ping()
	})
}
