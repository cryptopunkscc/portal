package main

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/getlantern/systray"
)

func newRunner(client portal.Client) *runner {
	return &runner{portal: client}
}

type runner struct {
	portal portal.Client
	log    plog.Logger
}

func (t *runner) Run(ctx context.Context) (err error) {
	if err = t.portal.Ping(); err != nil {
		return errors.New("portal-tray requires portal-app running")
	}

	t.log = plog.Get(ctx).Type(t).Set(&ctx)

	go func() {
		t.portal.Join()
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
					if err := t.portal.Open(nil, "launcher"); err != nil {
						t.log.Println("launcher:", err)
					}
				}()
			case <-quit.ClickedCh:
				if err := t.portal.Close(); err != nil {
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
