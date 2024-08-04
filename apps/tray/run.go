package main

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/getlantern/systray"
)

func newRunner(api portalApi) *runner {
	return &runner{api: api}
}

type runner struct {
	api portalApi
	log plog.Logger
}

func (t *runner) Run(ctx context.Context) (err error) {
	if err = t.api.Ping(); err != nil {
		return errors.New("portal-tray requires portal-app running")
	}

	t.log = plog.Get(ctx).Type(t).Set(&ctx)

	go func() {
		t.api.Await()
		systray.Quit()
	}()
	go func() {
		<-ctx.Done()
		systray.Quit()
	}()

	systray.Run(t.onReady, t.onExit)
	return
}

func (t *runner) onReady() {
	t.log.Println("portal tray start")
	launcher := systray.AddMenuItem("Launcher", "Launcher")
	quit := systray.AddMenuItem("Quit ", "Quit")

	go func() {
		for {
			select {
			case <-launcher.ClickedCh:
				go func() {
					if err := t.api.Open("launcher"); err != nil {
						t.log.Println("launcher:", err)
					}
				}()
			case <-quit.ClickedCh:
				if err := t.api.Close(); err != nil {
					t.log.Println("quit:", err)
					systray.Quit()
				}
			}
		}
	}()
}

func (t *runner) onExit() {
	t.log.Println("portal tray exit")
}
