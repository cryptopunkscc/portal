package main

import (
	"context"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os/exec"
	"time"
)

func startPortald(ctx context.Context, client apphost.Portald) (err error) {
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

func awaitPortaldService(ctx context.Context, client apphost.Portald) (err error) {
	await := flow.Await{
		UpTo:  5 * time.Second,
		Delay: 50 * time.Millisecond,
		Mod:   6,
		Ctx:   ctx,
	}
	for range await.Chan() {
		if err = apphost.Default.Connect(); err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	for range await.Chan() {
		if err = client.Ping(); err == nil {
			break
		}
	}
	return
}
