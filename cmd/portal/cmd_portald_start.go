package main

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os/exec"
	"time"
)

func (a *Application) startPortald(ctx context.Context) (err error) {
	if err = a.startPortaldProcess(ctx); err != nil {
		return
	}
	if err = a.awaitPortaldService(ctx); err != nil {
		return
	}
	return
}

func (a *Application) startPortaldProcess(ctx context.Context) (err error) {
	plog.Get(ctx).Println("starting portald")
	c := exec.Command("portald")
	err = c.Start()
	return
}

func (a *Application) awaitPortaldService(ctx context.Context) (err error) {
	await := flow.Await{
		UpTo:  5 * time.Second,
		Delay: 50 * time.Millisecond,
		Mod:   6,
		Ctx:   ctx,
	}
	for range await.Chan() {
		if err = a.Apphost.Connect(); err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	for range await.Chan() {
		if err = a.portald().Ping(); err == nil {
			break
		}
	}
	return
}
